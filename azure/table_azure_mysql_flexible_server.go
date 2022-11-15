package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2021-05-01/mysqlflexibleservers"
)

//// TABLE DEFINITION

func tableAzureMySQLFlexibleServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mysql_flexible_server",
		Description: "Azure MySQL Flexible Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMySQLFlexibleServer,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listMySQLFlexibleServers,
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
				Transform:   transform.FromField("ServerProperties.State"),
			},
			{
				Name:        "version",
				Description: "Specifies the version of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Version").Transform(transform.ToString),
			},
			{
				Name:        "administrator_login",
				Description: "The administrator's login name of a server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AdministratorLogin"),
			},
			{
				Name:        "availability_zone",
				Description: "Availability Zone information of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AvailabilityZone"),
			},
			{
				Name:        "backup_retention_days",
				Description: "Backup retention days for the server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.Backup.BackupRetentionDays"),
			},
			{
				Name:        "create_mode",
				Description: "The mode to create a new server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.CreateMode"),
			},
			{
				Name:        "earliest_restore_date",
				Description: "Specifies the earliest restore point creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServerProperties.Backup.EarliestRestoreDate").Transform(convertDateToTime),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.FullyQualifiedDomainName"),
			},
			{
				Name:        "geo_redundant_backup",
				Description: "Indicates whether Geo-redundant is enabled, or not for server backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Backup.GeoRedundantBackup").Transform(transform.ToString),
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
				Transform:   transform.FromField("ServerProperties.Network.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "replica_capacity",
				Description: "The maximum number of replicas that a primary server can have.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.ReplicaCapacity"),
			},
			{
				Name:        "replication_role",
				Description: "The replication role of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.ReplicationRole"),
			},
			{
				Name:        "restore_point_in_time",
				Description: "Restore point creation time (ISO8601 format), specifying the time to restore from.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServerProperties.RestorePointInTime"),
			},
			{
				Name:        "sku_name",
				Description: "The name of the sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "The tier of the particular SKU.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "source_server_resource_id",
				Description: "The source MySQL server id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.SourceServerResourceID"),
			},
			{
				Name:        "storage_auto_grow",
				Description: "Indicates whether storage auto grow is enabled, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Storage.AutoGrow").Transform(transform.ToString),
			},
			{
				Name:        "storage_iops",
				Description: "Storage IOPS for a server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.Storage.Iops"),
			},
			{
				Name:        "storage_size_gb",
				Description: "Indicates max storage allowed for a server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.Storage.StorageSizeGB"),
			},
			{
				Name:        "storage_sku",
				Description: "The sku name of the server storage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Storage.StorageSku"),
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
				Transform:   transform.FromField("ServerProperties.HighAvailability"),
			},
			{
				Name:        "network",
				Description: "Network related properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServerProperties.Network"),
			},
			{
				Name:        "maintenance_window",
				Description: "Maintenance window of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServerProperties.MaintenanceWindow"),
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
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysqlflexibleservers.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	resourceGroupName := h.Item.(resources.Group).Name

	result, err := client.ListByResourceGroup(ctx, *resourceGroupName)
	if err != nil {
		plugin.Logger(ctx).Error("listMySQLFlexibleServers", "list", err)
		return nil, err
	}

	for _, server := range result.Values() {
		d.StreamListItem(ctx, server)
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
		for _, server := range result.Values() {
			d.StreamListItem(ctx, server)
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

func getMySQLFlexibleServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMySQLFlexibleServer")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// check if name or resourceGroup is empty
	if resourceGroup == "" || name == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysqlflexibleservers.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getMySQLFlexibleServer", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func listMySQLFlexibleServersConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMySQLFlexibleServersConfigurations")

	server := h.Item.(mysqlflexibleservers.Server)
	resourceGroup := strings.Split(string(*server.ID), "/")[4]
	serverName := *server.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysqlflexibleservers.NewConfigurationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroup, serverName)
	if err != nil {
		plugin.Logger(ctx).Error("listMySQLFlexibleServersConfigurations", "list", err)
		return nil, err
	}

	var mySQLFlexibleServersConfigurations []map[string]interface{}

	for _, i := range op.Values() {
		mySQLFlexibleServersConfigurations = append(mySQLFlexibleServersConfigurations, extractMySQLFlexibleServersconfiguration(i))
	}

	return mySQLFlexibleServersConfigurations, nil
}

//// TRANSFORM FUNCTION

// If we return the API response directly, the output will not provide the properties of Configurations
func extractMySQLFlexibleServersconfiguration(i mysqlflexibleservers.Configuration) map[string]interface{} {
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
	if i.ConfigurationProperties != nil {
		mySQLFlexibleServersconfiguration["ConfigurationProperties"] = *i.ConfigurationProperties
	}

	return mySQLFlexibleServersconfiguration
}
