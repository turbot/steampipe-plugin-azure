package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	keyVaultp1 "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureKeyVaultCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault_certificate",
		Description: "Azure Key Vault Certificate",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"vault_name", "name"}),
			Hydrate:    getKeyVaultCertificate,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "404", "SecretDisabled"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate:       listKeyVaultCertificates,
			ParentHydrate: listKeyVaults,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name: "vault_name", Require: plugin.Optional,
				},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCertificateNameAndVaultName, "Name"),
			},
			{
				Name:        "vault_name",
				Description: "The name of the vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCertificateNameAndVaultName, "VaultName"),
			},
			// We are getting the ID value from Get API call correctly not from List API call
			// Get Response: https://turbottest94388.vault.azure.net/certificates/turbottest94388/beaf55112a214cd88aa500fcee10b0f4
			// List Response: https://turbottest94388.vault.azure.net/certificates/turbottest94388
			{
				Name:        "id",
				Description: "Certificate identifier.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "x509_thumbprint",
				Description: "Thumbprint of the certificate. A URL-encoded base64 string.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("X509Thumbprint"),
			},
			{
				Name:        "recovery_level",
				Description: "Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the certificate can be permanently deleted by a privileged user; otherwise, only the system can purge the certificate, at the end of the retention interval. Possible values include: 'Purgeable', 'RecoverablePurgeable', 'Recoverable', 'RecoverableProtectedSubscription'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Attributes.RecoveryLevel"),
			},
			{
				Name:        "enabled",
				Description: "Determines whether the object is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Attributes.Enabled"),
			},
			{
				Name:        "not_before",
				Description: "Not before date in UTC.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromP(convertPointerUnixTimestampToTimestamp, "NotBefore").Transform(transform.UnixMsToTimestamp).NullIfZero(),
			},
			{
				Name:        "expires",
				Description: "Expiry date in UTC.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromP(convertPointerUnixTimestampToTimestamp, "Expires").Transform(transform.UnixMsToTimestamp).NullIfZero(),
			},
			{
				Name:        "created",
				Description: "Creation time in UTC.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromP(convertPointerUnixTimestampToTimestamp, "Created").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "updated",
				Description: "Last updated time in UTC.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromP(convertPointerUnixTimestampToTimestamp, "Updated").Transform(transform.UnixMsToTimestamp).NullIfZero(),
			},
			{
				Name:        "key_id",
				Description: "The key id.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Kid"),
			},
			{
				Name:        "secret_id",
				Description: "The secret id.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Sid"),
			},
			{
				Name:        "content_type",
				Description: "The content type of the secret.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultCertificate,
			},
			{
				Name:        "cer",
				Description: "CER contents of x509 certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultCertificate,
			},
			{
				Name:        "key_properties",
				Description: "Properties of the key backing a certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Policy.KeyProperties"),
			},
			{
				Name:        "secret_properties",
				Description: "Properties of the secret backing a certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Policy.SecretProperties"),
			},
			{
				Name:        "x509_certificate_properties",
				Description: "Properties of the X509 component of a certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Policy.X509CertificateProperties"),
			},
			{
				Name:        "lifetime_actions",
				Description: "Actions that will be performed by Key Vault over the lifetime of a certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Policy.LifetimeActions"),
			},
			{
				Name:        "issuer_parameters",
				Description: "Parameters for the issuer of the X509 component of a certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultCertificate,
				Transform:   transform.FromField("Policy.IssuerParameters"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCertificateNameAndVaultName, "Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Hydrate:     getKeyVaultCertificate,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// We will not get the location and resource groups for the certificate because they are based on vault.
			// Azure standard columns
		}),
	}
}

//// LIST FUNCTION

func listKeyVaultCertificates(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of vault
	vault := h.Item.(keyVaultp1.Resource)

	vaultName := d.EqualsQualString("vault_name")

	if vaultName != "" && vaultName != *vault.Name{
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "VAULT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_certificate.listKeyVaultCertificates", "session_error", err)
		return nil, err
	}
	vaultURI := "https://" + *vault.Name + ".vault.azure.net/"

	client := keyvault.New()
	client.Authorizer = session.Authorizer

	maxresult := int32(25)

	result, err := client.GetCertificates(ctx, vaultURI, &maxresult)
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_certificate.listKeyVaultCertificates", "api_error", err)
		return nil, err
	}
	for _, cert := range result.Values() {
		d.StreamListItem(ctx, cert)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, cert := range result.Values() {
			d.StreamListItem(ctx, cert)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKeyVaultCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var vaultName, name string
	if h.Item != nil {
		data := h.Item.(keyvault.CertificateItem)
		splitID := strings.Split(*data.ID, "/")
		vaultName = strings.Split(splitID[2], ".")[0]
		name = splitID[4]

		// Operation get is not allowed on a disabled certificate
		if !*data.Attributes.Enabled {
			return nil, nil
		}
	} else {
		vaultName = d.EqualsQuals["vault_name"].GetStringValue()
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Create session
	session, err := GetNewSession(ctx, d, "VAULT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_certificate.getKeyVaultCertificate", "session_error", err)
		return nil, err
	}

	client := keyvault.New()
	client.Authorizer = session.Authorizer

	vaultURI := "https://" + vaultName + ".vault.azure.net/"

	op, err := client.GetCertificate(ctx, vaultURI, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_certificate.getKeyVaultCertificate", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func convertPointerUnixTimestampToTimestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	result := make(map[string]interface{}, 0)
	switch item := d.HydrateItem.(type) {
	case keyvault.CertificateItem:
		a := item.Attributes
		if a != nil {
			if a.Created != nil {
				result["Created"] = a.Created.Duration()
			}
			if a.Expires != nil {
				result["Expires"] = a.Expires.Duration()
			}
			if a.NotBefore != nil {
				result["NotBefore"] = a.NotBefore.Duration()
			}
			if a.Updated != nil {
				result["Updated"] = a.Updated.Duration()
			}
		}
	case keyvault.CertificateBundle:
		a := item.Attributes
		if a != nil {
			if a.Created != nil {
				result["Created"] = a.Created.Duration()
			}
			if a.Expires != nil {
				result["Expires"] = a.Expires.Duration()
			}
			if a.NotBefore != nil {
				result["NotBefore"] = a.NotBefore.Duration()
			}
			if a.Updated != nil {
				result["Updated"] = a.Updated.Duration()
			}
		}
	}

	return result[param], nil
}

func getCertificateNameAndVaultName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	result := make(map[string]interface{}, 0)
	if d.HydrateItem != nil {
		switch item := d.HydrateItem.(type) {
		case keyvault.CertificateItem:
			result["Name"] = strings.Split(*item.ID, "/")[4]
			result["VaultName"] = strings.Split(result["Name"].(string), ".")[0]
		case keyvault.CertificateBundle:
			result["Name"] = strings.Split(*item.ID, "/")[4]
			result["VaultName"] = strings.Split(result["Name"].(string), ".")[0]
		}
	}
	return result[param], nil
}
