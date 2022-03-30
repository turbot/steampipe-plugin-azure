package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	secret "github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureKeyVaultSecret(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault_secret",
		Description: "Azure Key Vault Secret",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"vault_name", "name"}),
			Hydrate:           getKeyVaultSecret,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "404", "SecretDisabled"}),
		},
		List: &plugin.ListConfig{
			Hydrate:       listKeyVaultSecrets,
			ParentHydrate: listKeyVaults,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractVaultNameFromSecretID, "Name"),
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a secret uniquely.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultSecret,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "vault_name",
				Description: "The friendly name that identifies the vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractVaultNameFromSecretID, "VaultName"),
			},
			{
				Name:        "enabled",
				Description: "Indicates whether the secret is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Attributes.Enabled"),
			},
			{
				Name:        "managed",
				Description: "Indicates whether the secret's lifetime is managed by key vault, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "content_type",
				Description: "Specifies the type of the secret value such as a password.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "Specifies the time when the secret is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Attributes.Created").Transform(convertDateUnixToTime),
			},
			{
				Name:        "expires_at",
				Description: "Specifies the time when the secret will expire.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Attributes.Expires").Transform(convertDateUnixToTime).Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "kid",
				Description: "If this is a secret backing a KV certificate, then this field specifies the corresponding key backing the KV certificate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultSecret,
			},
			{
				Name:        "not_before",
				Description: "Specifies the time before which the secret is not usable.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Attributes.NotBefore").Transform(convertDateUnixToTime),
			},
			{
				Name:        "recoverable_days",
				Description: "Specifies the soft delete data retention days. Value should be >=7 and <=90 when softDelete enabled, otherwise 0.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Attributes.RecoverableDays"),
			},
			{
				Name:        "recovery_level",
				Description: "The deletion recovery level currently in effect for the object. If it contains 'Purgeable', then the object can be permanently deleted by a privileged user; otherwise, only the system can purge the object at the end of the retention interval.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attributes.RecoveryLevel").Transform(transform.ToString),
			},
			{
				Name:        "updated_at",
				Description: "Specifies the time when the secret was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Attributes.Updated").Transform(convertDateUnixToTime),
			},
			{
				Name:        "value",
				Description: "Specifies the secret value.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultSecret,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractVaultNameFromSecretID, "Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTurbotData,
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTurbotData,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTurbotData,
			},
		}),
	}
}

//// LIST FUNCTION

func listKeyVaultSecrets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of key vault
	vault := h.Item.(keyvault.Resource)

	// Create session
	session, err := GetNewSession(ctx, d, "VAULT")
	if err != nil {
		return nil, err
	}

	vaultURI := "https://" + *vault.Name + ".vault.azure.net/"
	maxResults := int32(25)

	client := secret.New()
	client.Authorizer = session.Authorizer
	result, err := client.GetSecrets(ctx, vaultURI, &maxResults)
	if err != nil {
		/*
		* To make the above API call user must have list secrets permission.
		* If an user does not have this permission to perform list operation for secrets then the api returns forbiden error.
		* The permission should be grant in the key valult access policy.
		* If the forbidden error is not being handled here then it will make the table fail,
		  if there is at least a single key vault on which user does not have access to perform list secret operation.
		*/

		if strings.Contains(err.Error(), "Invalid audience. Expected https://vault.azure.net") {
			return nil, nil
		}
		return nil, err
	}

	for _, secret := range result.Values() {
		d.StreamLeafListItem(ctx, secret)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, secret := range result.Values() {
			d.StreamLeafListItem(ctx, secret)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKeyVaultSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVaultSecret")

	var vaultName, name string
	if h.Item != nil {
		data := h.Item.(secret.SecretItem)
		splitID := strings.Split(*data.ID, "/")
		vaultName = strings.Split(splitID[2], ".")[0]
		name = splitID[4]

		// Operation get is not allowed on a disabled secret
		if !*data.Attributes.Enabled {
			return nil, nil
		}
	} else {
		vaultName = d.KeyColumnQuals["vault_name"].GetStringValue()
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create session
	session, err := GetNewSession(ctx, d, "VAULT")
	if err != nil {
		return nil, err
	}

	client := secret.New()
	client.Authorizer = session.Authorizer

	vaultURI := "https://" + vaultName + ".vault.azure.net/"

	op, err := client.GetSecret(ctx, vaultURI, name, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func getTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getTurbotData")

	secretID := keyVaultSecretData(h.Item)
	splitID := strings.Split(secretID, "/")
	vaultName := strings.Split(splitID[2], ".")[0]

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := keyvault.NewVaultsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	maxResults := int32(100)

	op, err := client.List(ctx, &maxResults)
	if err != nil {
		return nil, err
	}

	var vaultID, location string
	for _, i := range op.Values() {
		if *i.Name == vaultName {
			vaultID = *i.ID
			location = *i.Location
		}
	}
	splitVaultID := strings.Split(vaultID, "/")
	akas := []string{"azure:///subscriptions/" + subscriptionID + "/resourceGroups/" + splitVaultID[4] + "/providers/Microsoft.KeyVault/vaults/" + vaultName + "/secrets/" + splitID[4], "azure:///subscriptions/" + subscriptionID + "/resourcegroups/" + splitVaultID[4] + "/providers/microsoft.keyvault/vaults/" + vaultName + "/secrets/" + splitID[4]}

	turbotData := map[string]interface{}{
		"SubscriptionId": subscriptionID,
		"ResourceGroup":  splitVaultID[4],
		"Location":       location,
		"Akas":           akas,
	}

	return turbotData, nil
}

//// TRANSFORM FUNCTIONS

func extractVaultNameFromSecretID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	secretID := keyVaultSecretData(d.HydrateItem)
	param := d.Param.(string)

	splitID := strings.Split(secretID, "/")

	result := map[string]string{
		"VaultName": strings.Split(splitID[2], ".")[0],
		"Name":      splitID[4],
	}

	return result[param], nil
}

func keyVaultSecretData(item interface{}) string {
	switch item := item.(type) {
	case secret.SecretItem:
		return *item.ID
	case secret.SecretBundle:
		return *item.ID
	}
	return ""
}
