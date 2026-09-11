package azure

import (
	"context"

	armredisenterprise "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v4"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureRedisEnterpriseCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_redis_enterprise_cluster",
		Description: "Azure Redis Enterprise Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getRedisEnterpriseCluster,
			Tags: map[string]string{
				"service": "Microsoft.Cache",
				"action":  "redisEnterprise/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listRedisEnterpriseClusters,
			Tags: map[string]string{
				"service": "Microsoft.Cache",
				"action":  "redisEnterprise/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Distinguishes the kind of cluster. Possible values: Redis, RedisEnterprise.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Kind").Transform(ptrToString),
			},
			{
				Name:        "provisioning_state",
				Description: "Current provisioning status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState").Transform(ptrToString),
			},
			{
				Name:        "resource_state",
				Description: "Current resource status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ResourceState").Transform(ptrToString),
			},
			{
				Name:        "redis_version",
				Description: "Version of redis the cluster supports, e.g. '6'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RedisVersion"),
			},
			{
				Name:        "host_name",
				Description: "DNS name of the cluster endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HostName"),
			},
			{
				Name:        "minimum_tls_version",
				Description: "The minimum TLS version for the cluster to support, e.g. '1.2'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinimumTLSVersion").Transform(ptrToString),
			},
			{
				Name:        "high_availability",
				Description: "Enabled by default. If disabled, the data set is not replicated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HighAvailability").Transform(ptrToString),
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public network traffic can access the cluster. Possible values: Enabled, Disabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PublicNetworkAccess").Transform(ptrToString),
			},
			{
				Name:        "redundancy_mode",
				Description: "Explains the current redundancy strategy of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RedundancyMode").Transform(ptrToString),
			},
			{
				Name:        "sku_name",
				Description: "The type of RedisEnterprise cluster to deploy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SKU.Name").Transform(ptrToString),
			},
			{
				Name:        "sku_capacity",
				Description: "The size of the RedisEnterprise cluster.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SKU.Capacity"),
			},
			{
				Name:        "zones",
				Description: "The Availability Zones where this cluster will be deployed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "The identity of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "encryption",
				Description: "Encryption-at-rest configuration for the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Encryption"),
			},
			{
				Name:        "maintenance_configuration",
				Description: "Cluster-level maintenance configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.MaintenanceConfiguration"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "List of private endpoint connections associated with the specified RedisEnterprise cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PrivateEndpointConnections"),
			},
			{
				Name:        "system_data",
				Description: "Azure Resource Manager metadata containing createdBy and modifiedBy information.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SystemData"),
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
				Transform:   transform.FromField("Location").Transform(formatRegion).Transform(toLower),
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

func listRedisEnterpriseClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRedisEnterpriseClusters")

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_cluster.listRedisEnterpriseClusters", "session_error", err)
		return nil, err
	}

	d.WaitForListRateLimit(ctx)

	client, err := armredisenterprise.NewClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_cluster.listRedisEnterpriseClusters", "client_error", err)
		return nil, err
	}

	pager := client.NewListPager(nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_redis_enterprise_cluster.listRedisEnterpriseClusters", "api_error", err)
			return nil, err
		}
		for _, cluster := range result.Value {
			d.StreamListItem(ctx, *cluster)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedisEnterpriseCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRedisEnterpriseCluster")

	name := d.EqualsQualString("name")
	resourceGroup := d.EqualsQualString("resource_group")

	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_cluster.getRedisEnterpriseCluster", "session_error", err)
		return nil, err
	}

	client, err := armredisenterprise.NewClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_cluster.getRedisEnterpriseCluster", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_redis_enterprise_cluster.getRedisEnterpriseCluster", "api_error", err)
		return nil, err
	}

	if op.ID != nil {
		return op.Cluster, nil
	}

	return nil, nil
}
