package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2020-01-01/mariadb"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureMariaDBServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mariadb_server",
		Description: "Azure MariaDB Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMariaDBServer,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "400", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMariaDBServers,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "A fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "Specifies the server version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Version").Transform(transform.ToString),
			},
			{
				Name:        "geo_redundant_backup_enabled",
				Description: "Indicates whether geo-redundant backup is enabled for server backup, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.StorageProfile.GeoRedundantBackup").Transform(transform.ToString),
			},
			{
				Name:        "user_visible_state",
				Description: "A state of a server that is visible to user. Valid values are: 'Ready', 'Dropping', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.UserVisibleState").Transform(transform.ToString),
			},
			{
				Name:        "administrator_login",
				Description: "The administrator's login name of a server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AdministratorLogin"),
			},
			{
				Name:        "auto_grow_enabled",
				Description: "Indicates whether storage auto grow is enabled for server, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.StorageProfile.StorageAutogrow").Transform(transform.ToString),
			},
			{
				Name:        "backup_retention_days",
				Description: "Specifies the backup retention days for the server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.StorageProfile.BackupRetentionDays"),
			},
			{
				Name:        "earliest_restore_date",
				Description: "Specifies the earliest restore point creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServerProperties.EarliestRestoreDate").Transform(convertDateToTime),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of a server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.FullyQualifiedDomainName"),
			},
			{
				Name:        "master_service_id",
				Description: "The master server id of a replica server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.MasterServerID"),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether or not public network access is allowed for this server. Valid values are: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "replication_role",
				Description: "The replication role of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.ReplicationRole"),
			},
			{
				Name:        "replica_capacity",
				Description: "The maximum number of replicas that a master server can have.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.ReplicaCapacity"),
			},
			{
				Name:        "sku_capacity",
				Description: "The scale up/out capacity, representing server's compute units.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_family",
				Description: "The family of hardware.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Family"),
			},
			{
				Name:        "sku_name",
				Description: "The name of the sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_size",
				Description: "The size code, to be interpreted by resource as appropriate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Size"),
			},
			{
				Name:        "sku_tier",
				Description: "The tier of the particular SKU. Valid values are: 'Basic', 'GeneralPurpose', 'MemoryOptimized'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "ssl_enforcement",
				Description: "Indicates whether SSL enforcement is enabled, or not. Valid values are: 'Enabled', and 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.SslEnforcement").Transform(transform.ToString),
			},
			{
				Name:        "storage_mb",
				Description: "Specifies the max storage allowed for a server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.StorageProfile.StorageMB"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "A list of private endpoint connections on a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServerProperties.PrivateEndpointConnections"),
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

func listMariaDBServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMariaDBServers")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := mariadb.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, server := range *result.Value {
		d.StreamListItem(ctx, server)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMariaDBServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMariaDBServer")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := mariadb.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
