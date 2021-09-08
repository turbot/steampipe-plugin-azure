package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
)

//// TABLE DEFINITION

func tableAzureMSSQLVirtualMachine(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mssql_virtual_machine",
		Description: "Azure MS SQL Virtual Machine",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getMSSQLVirtualMachine,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listMSSQLVirtualMachines,
		},
		Columns: []*plugin.Column{
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
				Description: "SQL image offer.",
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
				Description: "SQL Server license type. Possible values include: 'PAYG', 'AHUB', 'DR'.",
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
				Description: "Azure Active Directory identity of the server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "key_vault_credential_settings",
				Description: "Key vault credential settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.KeyVaultCredentialSettings"),
			},
			{
				Name:        "server_configurations_management_settings",
				Description: "SQL Server configuration management settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ServerConfigurationsManagementSettings"),
			},
			{
				Name:        "storage_configuration_settings",
				Description: "Storage Configuration Settings.",
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listMSSQLVirtualMachines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sqlvirtualmachine.NewSQLVirtualMachinesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}
	
	for _, virtualMachine := range result.Values() {
		d.StreamListItem(ctx, virtualMachine)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, virtualMachine := range result.Values() {
			d.StreamListItem(ctx, virtualMachine)
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getMSSQLVirtualMachine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMSSQLVirtualMachine")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sqlvirtualmachine.NewSQLVirtualMachinesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name, "")
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
