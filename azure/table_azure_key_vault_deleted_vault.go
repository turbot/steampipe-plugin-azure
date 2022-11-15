package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/mgmt/keyvault"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAzureKeyVaultDeletedVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault_deleted_vault",
		Description: "Azure Key Vault Deleted Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "region"}),
			Hydrate:    getKeyVaultDeletedVault,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyVaultDeletedVaults,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the deleted vault.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a deleted vault uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "vault_id",
				Description: "The resource id of the original vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.VaultID"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_date",
				Description: "The deleted date of the vault.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.DeletionDate").Transform(convertDateToTime),
			},
			{
				Name:        "scheduled_purge_date",
				Description: "The scheduled purged date of the vault.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.ScheduledPurgeDate").Transform(convertDateToTime),
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
				Transform:   transform.FromField("Properties.Tags"),
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
				Transform:   transform.FromField("Properties.Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.VaultID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listKeyVaultDeletedVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	keyVaultClient := keyvault.NewVaultsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	keyVaultClient.Authorizer = session.Authorizer

	result, err := keyVaultClient.ListDeleted(ctx)
	if err != nil {
		return nil, err
	}
	for _, deletedVault := range result.Values() {
		d.StreamListItem(ctx, deletedVault)
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
		for _, deletedVault := range result.Values() {
			d.StreamListItem(ctx, deletedVault)
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

func getKeyVaultDeletedVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVaultDeletedVault")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	name := d.KeyColumnQuals["name"].GetStringValue()
	region := d.KeyColumnQuals["region"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || region == "" {
		return nil, nil
	}

	client := keyvault.NewVaultsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.GetDeleted(ctx, name, region)
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
