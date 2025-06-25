package azure

import (
	"context"
	"slices"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	armmysqlflexibleservers "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysqlflexibleservers/v2"
)

//// TABLE DEFINITION

func tableAzureMySQLFlexibleServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mysql_flexible_server",
		Description: "Azure MySQL Flexible Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMySQLFlexibleServer,
			Tags: map[string]string{
				"service": "Microsoft.DBforMySQL",
				"action":  "flexibleServers/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listMySQLFlexibleServers,
			Tags: map[string]string{
				"service": "Microsoft.DBforMySQL",
				"action":  "flexibleServers/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a server uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.State"),
			},
			{
				Name:        "version",
				Description: "Specifies the version of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Version"),
			},
			{
				Name:        "administrator_login",
				Description: "The administrator's login name of a server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AdministratorLogin"),
			},
			{
				Name:        "availability_zone",
				Description: "Availability Zone information of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AvailabilityZone"),
			},
			{
				Name:        "backup_retention_days",
				Description: "Backup retention days for the server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Backup.BackupRetentionDays"),
			},
			{
				Name:        "create_mode",
				Description: "The mode to create a new server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CreateMode"),
			},
			{
				Name:        "earliest_restore_date",
				Description: "Specifies the earliest restore point creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.Backup.EarliestRestoreDate"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.FullyQualifiedDomainName"),
			},
			{
				Name:        "geo_redundant_backup",
				Description: "Indicates whether Geo-redundant is enabled, or not for server backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Backup.GeoRedundantBackup"),
			},
			{
				Name:        "location",
				Description: "The server location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public network access is allowed for this server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Network.PublicNetworkAccess"),
			},
			{
				Name:        "replica_capacity",
				Description: "The maximum number of replicas that a primary server can have.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.ReplicaCapacity"),
			},
			{
				Name:        "replication_role",
				Description: "The replication role of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ReplicationRole"),
			},
			{
				Name:        "restore_point_in_time",
				Description: "Restore point creation time (ISO8601 format), specifying the time to restore from.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.RestorePointInTime"),
			},
			{
				Name:        "sku_name",
				Description: "The name of the sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SKU.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "The tier of the particular SKU.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SKU.Tier"),
			},
			{
				Name:        "source_server_resource_id",
				Description: "The source MySQL server id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SourceServerResourceID"),
			},
			{
				Name:        "storage_auto_grow",
				Description: "Indicates whether storage auto grow is enabled, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Storage.AutoGrow"),
			},
			{
				Name:        "storage_iops",
				Description: "Storage IOPS for a server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Storage.Iops"),
			},
			{
				Name:        "storage_size_gb",
				Description: "Indicates max storage allowed for a server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Storage.StorageSizeGB"),
			},
			{
				Name:        "storage_sku",
				Description: "The sku name of the server storage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Storage.StorageSKU"),
			},
			{
				Name:        "flexible_server_configurations",
				Description: "The server configurations(parameters) details of the server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMySQLFlexibleServersConfigurations,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "high_availability",
				Description: "High availability related properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.HighAvailability"),
			},
			{
				Name:        "network",
				Description: "Network related properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Network"),
			},
			{
				Name:        "maintenance_window",
				Description: "Maintenance window of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.MaintenanceWindow"),
			},
			{
				Name:        "system_data",
				Description: "The system metadata relating to this server.",
				Type:        proto.ColumnType_JSON,
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

func listMySQLFlexibleServers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := armmysqlflexibleservers.NewServersClient(subscriptionID, session.Cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_flexible_server.listMySQLFlexibleServers", "session_error", err)
		return nil, err
	}
	// client.Authorizer = session.Authorizer

	resourceGroupName := h.Item.(resources.Group).Name

	pager := client.NewListByResourceGroupPager(*resourceGroupName, nil)

	for pager.More() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mysql_flexible_server.listMySQLFlexibleServers", "api_error", err)
			return nil, err
		}

		for _, item := range page.Value {
			d.StreamListItem(ctx, *item)
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

func getMySQLFlexibleServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// check if name or resourceGroup is empty
	if resourceGroup == "" || name == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_flexible_server.getMySQLFlexibleServer", "session_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID

	client, err := armmysqlflexibleservers.NewServersClient(subscriptionID, session.Cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_flexible_server.getMySQLFlexibleServer", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_flexible_server.getMySQLFlexibleServer", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op.Server, nil
	}

	return nil, nil
}

func listMySQLFlexibleServersConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	server := h.Item.(armmysqlflexibleservers.Server)
	resourceGroup := strings.Split(string(*server.ID), "/")[4]
	serverName := *server.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_flexible_server.listMySQLFlexibleServersConfigurations", "session_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID

	client, err := armmysqlflexibleservers.NewConfigurationsClient(subscriptionID, session.Cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_flexible_server.listMySQLFlexibleServersConfigurations", "client_error", err)
		return nil, err
	}

	// Return nil for the specific states of the flexible server where the
	// API does not allow any operations to be performed.

	// Since it is challenging to test all possible states("Disabled", "Dropping", "Ready", "Starting", "Stopped", "Stopping", "Updating") of the flexible server,
	// we have added a check to restrict
	// API calls for the "Stopping" and "Stopped" states based on current testing.

	// This logic may need to be updated in the future if the API behavior changes or if additional states need to be handled.

	// 1. Stopping:
	// Error: azure: GET https://management.azure.com/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/new-rg/providers/Microsoft.DBforMySQL/flexibleServers/test53/configurations
	// --------------------------------------------------------------------------------
	// RESPONSE 503: 503 Service Unavailable
	// ERROR CODE: ServiceBusy
	// --------------------------------------------------------------------------------
	// {
	//   "error": {
	//     "code": "ServiceBusy",
	//     "message": "Service is temporarily busy and the operation cannot be performed. Please try again later."
	//   }
	// }
	// 2. Stopped
	// Error: azure: GET https://management.azure.com/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/new-rg/providers/Microsoft.DBforMySQL/flexibleServers/test53/configurations
	// --------------------------------------------------------------------------------
	// RESPONSE 409: 409 Conflict
	// ERROR CODE: ServerUnavailableForOperation
	// --------------------------------------------------------------------------------
	// {
	//   "error": {
	// 	"code": "ServerUnavailableForOperation",
	// 	"message": "Operation 'GetServerParameters' cannot be performed as server 'test53' is currently in state 'Stopped'."
	//   }
	// }
	// 3. Starting: We are not getting any error
	if server.Properties != nil && slices.Contains([]string{"Stopping", "Stopped", "Updating"}, string(*server.Properties.State)) {
		return nil, nil
	}
	pager := client.NewListByServerPager(resourceGroup, serverName, nil)

	var mySQLFlexibleServersConfigurations []map[string]interface{}

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mysql_flexible_server.listMySQLFlexibleServersConfigurations", "api_error", err)
			return nil, err
		}

		for _, item := range page.Value {
			mySQLFlexibleServersConfigurations = append(mySQLFlexibleServersConfigurations, extractMySQLFlexibleServersconfiguration(item))
		}
	}

	return mySQLFlexibleServersConfigurations, nil
}

//// TRANSFORM FUNCTION

// If we return the API response directly, the output will not provide the properties of Configurations
func extractMySQLFlexibleServersconfiguration(i *armmysqlflexibleservers.Configuration) map[string]interface{} {
	mySQLFlexibleServersconfiguration := make(map[string]interface{})

	if i.ID != nil {
		mySQLFlexibleServersconfiguration["ID"] = *i.ID
	}
	if i.Name != nil {
		mySQLFlexibleServersconfiguration["Name"] = *i.Name
	}
	if i.Type != nil {
		mySQLFlexibleServersconfiguration["Type"] = *i.Type
	}
	if i.Properties != nil {
		mySQLFlexibleServersconfiguration["ConfigurationProperties"] = *i.Properties
	}

	return mySQLFlexibleServersconfiguration
}
