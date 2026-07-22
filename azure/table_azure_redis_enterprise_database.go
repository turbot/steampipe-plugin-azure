package azure

import (
	"context"
	"strings"

	armredisenterprise "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v4"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type redisEnterpriseDatabaseInfo struct {
	Database      armredisenterprise.Database
	ClusterName   *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION

func tableAzureRedisEnterpriseDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_redis_enterprise_database",
		Description: "Azure Redis Enterprise Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"cluster_name", "name", "resource_group"}),
			Hydrate:    getRedisEnterpriseDatabase,
			Tags: map[string]string{
				"service": "Microsoft.Cache",
				"action":  "redisEnterprise/databases/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listRedisEnterpriseClusters,
			Hydrate:       listRedisEnterpriseDatabases,
			Tags: map[string]string{
				"service": "Microsoft.Cache",
				"action":  "redisEnterprise/databases/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the database resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Name"),
			},
			{
				Name:        "cluster_name",
				Description: "The name of the RedisEnterprise cluster this database belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterName"),
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.ID"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Type"),
			},
			{
				Name:        "provisioning_state",
				Description: "Current provisioning status of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.ProvisioningState").Transform(ptrToString),
			},
			{
				Name:        "resource_state",
				Description: "Current resource status of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.ResourceState").Transform(ptrToString),
			},
			{
				Name:        "redis_version",
				Description: "Version of Redis the database is running on, e.g. '6.0' or '7.4'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.RedisVersion"),
			},
			{
				Name:        "client_protocol",
				Description: "Specifies whether redis clients can connect using TLS-encrypted or plaintext redis protocols. Possible values: Encrypted, Plaintext.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.ClientProtocol").Transform(ptrToString),
			},
			{
				Name:        "port",
				Description: "TCP port of the database endpoint.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Database.Properties.Port"),
			},
			{
				Name:        "clustering_policy",
				Description: "Clustering policy. Possible values: EnterpriseCluster, OSSCluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.ClusteringPolicy").Transform(ptrToString),
			},
			{
				Name:        "eviction_policy",
				Description: "Redis eviction policy. Possible values: AllKeysLFU, AllKeysLRU, AllKeysRandom, VolatileLRU, VolatileLFU, VolatileTTL, VolatileRandom, NoEviction.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.EvictionPolicy").Transform(ptrToString),
			},
			{
				Name:        "access_keys_authentication",
				Description: "Whether access with access keys is enabled. Possible values: Enabled, Disabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.AccessKeysAuthentication").Transform(ptrToString),
			},
			{
				Name:        "defer_upgrade",
				Description: "Option to defer upgrade when newest version is released. Possible values: Deferred, NotDeferred.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Properties.DeferUpgrade").Transform(ptrToString),
			},
			{
				Name:        "persistence",
				Description: "Persistence settings (AOF/RDB configuration).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Database.Properties.Persistence"),
			},
			{
				Name:        "modules",
				Description: "Optional set of redis modules enabled in this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Database.Properties.Modules"),
			},
			{
				Name:        "geo_replication",
				Description: "Optional set of properties to configure geo replication for this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Database.Properties.GeoReplication"),
			},
			{
				Name:        "system_data",
				Description: "Azure Resource Manager metadata containing createdBy and modifiedBy information.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Database.SystemData"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Database.ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(formatRegion).Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceGroup"),
			},
		}),
	}
}

//// LIST FUNCTION

func listRedisEnterpriseDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRedisEnterpriseDatabases")

	cluster := h.Item.(armredisenterprise.Cluster)
	clusterName := cluster.Name
	resourceGroup := &strings.Split(*cluster.ID, "/")[4]

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_database.listRedisEnterpriseDatabases", "session_error", err)
		return nil, err
	}

	dbClient, err := armredisenterprise.NewDatabasesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_database.listRedisEnterpriseDatabases", "client_error", err)
		return nil, err
	}

	pager := dbClient.NewListByClusterPager(*resourceGroup, *clusterName, nil)
	for pager.More() {
		d.WaitForListRateLimit(ctx)
		result, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_redis_enterprise_database.listRedisEnterpriseDatabases", "api_error", err)
			return nil, err
		}
		for _, db := range result.Value {
			d.StreamLeafListItem(ctx, redisEnterpriseDatabaseInfo{*db, clusterName, resourceGroup, cluster.Location})
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedisEnterpriseDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRedisEnterpriseDatabase")

	name := d.EqualsQualString("name")
	clusterName := d.EqualsQualString("cluster_name")
	resourceGroup := d.EqualsQualString("resource_group")

	if name == "" || clusterName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_database.getRedisEnterpriseDatabase", "session_error", err)
		return nil, err
	}

	dbClient, err := armredisenterprise.NewDatabasesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_database.getRedisEnterpriseDatabase", "client_error", err)
		return nil, err
	}

	db, err := dbClient.Get(ctx, resourceGroup, clusterName, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_database.getRedisEnterpriseDatabase", "api_error", err)
		return nil, err
	}

	// Fetch the parent cluster to get Location
	clusterClient, err := armredisenterprise.NewClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_database.getRedisEnterpriseDatabase", "cluster_client_error", err)
		return nil, err
	}

	cluster, err := clusterClient.Get(ctx, resourceGroup, clusterName, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_database.getRedisEnterpriseDatabase", "cluster_api_error", err)
		return nil, err
	}

	if db.ID != nil {
		return redisEnterpriseDatabaseInfo{db.Database, &clusterName, &resourceGroup, cluster.Location}, nil
	}

	return nil, nil
}
