package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	sqlvirtualmachine "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sqlvirtualmachine/armsqlvirtualmachine"
)

//// TABLE DEFINITION

func tableAzureMSSQLVirtualMachine(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mssql_virtual_machine",
		Description: "Azure MS SQL Virtual Machine",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMSSQLVirtualMachine,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMSSQLVirtualMachines,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state to track the async operation status.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sql_image_offer",
				Description: "SQL image offer for the SQL virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SQLImageOffer"),
			},
			{
				Name:        "sql_image_sku",
				Description: "SQL Server edition type. Possible values include: 'Developer', 'Express', 'Standard', 'Enterprise', 'Web'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SQLImageSku"),
			},
			{
				Name:        "sql_management",
				Description: "SQL Server Management type. Possible values include: 'Full', 'LightWeight', 'NoAgent'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SQLManagement"),
			},
			{
				Name:        "sql_server_license_type",
				Description: "SQL server license type for the SQL virtual machine. Possible values include: 'PAYG', 'AHUB', 'DR'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SQLServerLicenseType"),
			},
			{
				Name:        "sql_virtual_machine_group_resource_id",
				Description: "ARM resource id of the SQL virtual machine group this SQL virtual machine is or will be part of.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SQLVirtualMachineGroupResourceID"),
			},
			{
				Name:        "virtual_machine_resource_id",
				Description: "ARM resource id of underlying virtual machine created from SQL marketplace image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.VirtualMachineResourceID"),
			},
			{
				Name:        "auto_backup_settings",
				Description: "Auto backup settings for SQL Server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AutoBackupSettings"),
			},
			{
				Name:        "auto_patching_settings",
				Description: "Auto patching settings for applying critical security updates to SQL virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AutoPatchingSettings"),
			},
			{
				Name:        "identity",
				Description: "Azure Active Directory identity for the SQL virtual machine.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "key_vault_credential_settings",
				Description: "Key vault credential settings for the SQL virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.KeyVaultCredentialSettings"),
			},
			{
				Name:        "server_configurations_management_settings",
				Description: "SQL server configuration management settings for the SQL virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ServerConfigurationsManagementSettings"),
			},
			{
				Name:        "storage_configuration_settings",
				Description: "Storage configuration settings for the SQL virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.StorageConfigurationSettings"),
			},
			{
				Name:        "wsfc_domain_credentials",
				Description: "Domain credentials for setting up Windows Server Failover Cluster for SQL availability group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.WsfcDomainCredentials"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
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

func listMSSQLVirtualMachines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_virtual_machine.listMSSQLVirtualMachines", "connection error", err)
	}
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_virtual_machine.listMSSQLVirtualMachines", "session error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sqlvirtualmachine.NewSQLVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_virtual_machine.listMSSQLVirtualMachines", "client error", err)
	}

	pager := client.NewListPager(nil)

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mssql_virtual_machine.listMSSQLVirtualMachines", "api error", err)
		}
		for _, virtualMachine := range nextResult.Value {
			d.StreamListItem(ctx, virtualMachine)
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

func getMSSQLVirtualMachine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_virtual_machine.getMSSQLVirtualMachine", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client, err := sqlvirtualmachine.NewSQLVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_virtual_machine.getMSSQLVirtualMachine", "client error", err)
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_virtual_machine.getMSSQLVirtualMachine", "api error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
