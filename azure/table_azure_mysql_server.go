package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
)

//// TABLE DEFINITION

func tableAzureMySQLServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mysql_server",
		Description: "Azure MySQL Server",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getMySQLServer,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listMySQLServers,
		},
		Columns: []*plugin.Column{
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
				Name:        "user_visible_state",
				Description: "A state of a server that is visible to user. Possible values include: 'Ready', 'Dropping', 'Disabled', 'Inaccessible'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.UserVisibleState").Transform(transform.ToString),
			},
			{
				Name:        "location",
				Description: "The resource location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrator_login",
				Description: "Specifies the username of the administrator for this server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AdministratorLogin"),
			},
			{
				Name:        "backup_retention_days",
				Description: "Backup retention days for the server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.StorageProfile.BackupRetentionDays"),
			},
			{
				Name:        "byok_enforcement",
				Description: "Status showing whether the server data encryption is enabled with customer-managed keys.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.ByokEnforcement"),
			},
			{
				Name:        "earliest_restore_date",
				Description: "Specifies the earliest restore point creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServerProperties.EarliestRestoreDate").Transform(convertDateToTime),
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
				Transform:   transform.FromField("ServerProperties.StorageProfile.GeoRedundantBackup").Transform(transform.ToString),
			},
			{
				Name:        "infrastructure_encryption",
				Description: "Status showing whether the server enabled infrastructure encryption. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.InfrastructureEncryption").Transform(transform.ToString),
			},
			{
				Name:        "master_server_id",
				Description: "The master server id of a replica server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.MasterServerID"),
			},
			{
				Name:        "minimal_tls_version",
				Description: "Enforce a minimal Tls version for the server. Possible values include: 'TLS10', 'TLS11', 'TLS12', 'TLSEnforcementDisabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.MinimalTLSVersion").Transform(transform.ToString),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether or not public network access is allowed for this server. Value is optional but if passed in, must be 'Enabled' or 'Disabled'. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "replica_capacity",
				Description: "The maximum number of replicas that a master server can have.",
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
				Description: "The name of the sku. For example: 'B_Gen4_1', 'GP_Gen5_8'.",
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
				Description: "The tier of the particular SKU. Possible values include: 'Basic', 'GeneralPurpose', 'MemoryOptimized'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "ssl_enforcement",
				Description: "Enable ssl enforcement or not when connect to server. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.SslEnforcement").Transform(transform.ToString),
			},
			{
				Name:        "storage_auto_grow",
				Description: "Indicates whether storage auto grow is enabled, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.StorageProfile.StorageAutogrow").Transform(transform.ToString),
			},
			{
				Name:        "storage_mb",
				Description: "Indicates max storage allowed for a server.",
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

func listMySQLServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysql.NewServersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	// Currently the API does not support pagination
	for _, server := range *result.Value {
		d.StreamListItem(ctx, server)
		// This will return zero if context has been cancelled (i.e due to manual cancellation) or
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMySQLServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMySQLServer")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Error: mysql.ServersClient#Get: Invalid input: autorest/validation: validation failed: parameter=resourceGroupName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysql.NewServersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
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
