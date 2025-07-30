package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDataProtectionBackupVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_protection_backup_vault",
		Description: "Azure Data Protection Backup Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAzureDataProtectionBackupVault,
			Tags: map[string]string{
				"service": "Microsoft.DataProtection",
				"action":  "backupVaults/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureDataProtectionBackupVaults,
			Tags: map[string]string{
				"service": "Microsoft.DataProtection",
				"action":  "backupVaults/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the backup vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the backup vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the backup vault resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "resource_move_state",
				Description: "The resource move state for the backup vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ResourceMoveState"),
			},
			{
				Name:        "resource_move_details",
				Description: "Details of the resource move for the backup vault.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ResourceMoveDetails"),
			},
			{
				Name:        "storage_settings",
				Description: "The storage settings of the backup vault.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.StorageSettings"),
			},
			{
				Name:        "monitoring_settings",
				Description: "The Monitoring Settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.MonitoringSettings"),
			},
			{
				Name:        "identity",
				Description: "Input Managed Identity Details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "system_data",
				Description: "Metadata pertaining to creation and last modification of the resource.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: "Tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: "The Azure region where the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: "The resource group in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listAzureDataProtectionBackupVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureDataProtectionBackupVaults")
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armdataprotection.NewBackupVaultsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_data_protection_backup_vault.listAzureDataProtectionBackupVaults", "client_error", err)
		return nil, err
	}

	input := &armdataprotection.BackupVaultsClientGetInSubscriptionOptions{}

	pager := clientFactory.NewGetInSubscriptionPager(input)
	if err != nil {
		plugin.Logger(ctx).Error("listAzureDataProtectionBackupVaults", "list_err", err)
		return nil, err
	}

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_data_protection_backup_vault.listAzureDataProtectionBackupVaults", "api_error", err)
			return nil, err
		}
		for _, backupVault := range page.Value {
			d.StreamListItem(ctx, backupVault)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getAzureDataProtectionBackupVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armdataprotection.NewBackupVaultsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_data_protection_backup_vault.getAzureDataProtectionBackupVault", "client_error", err)
		return nil, err
	}

	backupVault, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_data_protection_backup_vault.getAzureDataProtectionBackupVault", "api_error", err)
		return nil, err
	}

	if backupVault.ID != nil {
		return backupVault, nil
	}

	return nil, nil
}
