package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureKeyVaultKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault_key",
		Description: "Azure Key Vault Key",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"vault_name", "name", "resource_group"}),
			Hydrate:           getKeyVaultKey,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate:       listKeyVaultKeys,
			ParentHydrate: listKeyVaults,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a key uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "vault_name",
				Description: "The friendly name that identifies the vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractVaultNameFromID),
			},
			{
				Name:        "enabled",
				Description: "Indicates whether the key is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("KeyProperties.Attributes.Enabled"),
			},
			{
				Name:        "key_type",
				Description: "The type of the key. Possible values are: 'EC', 'ECHSM', 'RSA', 'RSAHSM'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultKey,
				Transform:   transform.FromField("KeyProperties.Kty").Transform(transform.ToString),
			},
			{
				Name:        "created_at",
				Description: "Specifies the time when the key is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.Created").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "curve_name",
				Description: "The elliptic curve name. Possible values are: 'P256', 'P384', 'P521', 'P256K'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultKey,
				Transform:   transform.FromField("KeyProperties.CurveName").Transform(transform.ToString),
			},
			{
				Name:        "expires_at",
				Description: "Specifies the time when the key wil expire.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.Expires").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "key_size",
				Description: "The key size in bits.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKeyVaultKey,
				Transform:   transform.FromField("KeyProperties.KeySize"),
			},
			{
				Name:        "key_uri",
				Description: "The URI to retrieve the current version of the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyProperties.KeyURI"),
			},
			{
				Name:        "key_uri_with_version",
				Description: "The URI to retrieve the specific version of the key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultKey,
				Transform:   transform.FromField("KeyProperties.KeyURIWithVersion"),
			},
			{
				Name:        "location",
				Description: "Azure location of the key vault resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "not_before",
				Description: "Specifies the time before which the key is not usable.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.NotBefore").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "recovery_level",
				Description: "The deletion recovery level currently in effect for the object. If it contains 'Purgeable', then the object can be permanently deleted by a privileged user; otherwise, only the system can purge the object at the end of the retention interval.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyProperties.Attributes.RecoveryLevel").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "Specifies the time when the key was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.Updated").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "key_ops",
				Description: "A list of key operations.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultKey,
				Transform:   transform.FromField("KeyProperties.KeyOps"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listKeyVaultKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of key vault
	vault := h.Item.(keyvault.Resource)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*vault.ID, "/")[4]

	client := keyvault.NewKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.List(ctx, resourceGroup, *vault.Name)
	if err != nil {
		return nil, err
	}
	for _, key := range result.Values() {
		d.StreamListItem(ctx, key)
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

		for _, key := range result.Values() {
			d.StreamListItem(ctx, key)
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

func getKeyVaultKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVaultKey")

	var vaultName, name, resourceGroup string
	if h.Item != nil {
		data := h.Item.(keyvault.Key)
		splitID := strings.Split(*data.ID, "/")
		vaultName = splitID[8]
		name = *data.Name
		resourceGroup = splitID[4]
	} else {
		vaultName = d.KeyColumnQuals["vault_name"].GetStringValue()
		name = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroup = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := keyvault.NewKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, vaultName, name)
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

//// TRANSFORM FUNCTIONS

func extractVaultNameFromID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(keyvault.Key)
	vaultName := strings.Split(*data.ID, "/")[8]
	return vaultName, nil
}
