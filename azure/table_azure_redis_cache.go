package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2020-06-01/redis"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureRedisCache(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_redis_cache",
		Description: "Azure Redis Cache",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getRedisCache,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "400", "400"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listRedisCaches,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the redis instance at the time the operation was called. Valid values are: 'Creating', 'Deleting', 'Disabled', 'Failed', 'Linking', 'Provisioning', 'RecoveringScaleFailure', 'Scaling', 'Succeeded', 'Unlinking', 'Unprovisioning', and 'Updating'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "redis_version",
				Description: "Specifies the version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RedisVersion"),
			},
			{
				Name:        "enable_non_ssl_port",
				Description: "Specifies whether the non-ssl Redis server port (6379) is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnableNonSslPort"),
			},
			{
				Name:        "host_name",
				Description: "Specifies the name of the redis host.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HostName"),
			},
			{
				Name:        "minimum_tls_version",
				Description: "Specifies the TLS version requires to connect.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinimumTLSVersion").Transform(transform.ToString).Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "port",
				Description: "Specifies the redis non-SSL port.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Port"),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether or not public endpoint access is allowed for this cache. Valid values are: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "sku_capacity",
				Description: "The size of the Redis cache to deploy.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Sku.Capacity"),
			},
			{
				Name:        "sku_family",
				Description: "The SKU family to use.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Sku.Family").Transform(transform.ToString),
			},
			{
				Name:        "sku_name",
				Description: "The type of Redis cache to deploy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "ssl_port",
				Description: "Specifies the redis SSL port.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.SslPort"),
			},
			{
				Name:        "subnet_id",
				Description: "The full resource ID of a subnet in a virtual network to deploy the Redis cache in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SubnetID"),
			},
			{
				Name:        "static_ip",
				Description: "Specifies the statis IP address. Required when deploying a Redis cache inside an existing Azure Virtual Network.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Properties.StaticIP"),
			},
			{
				Name:        "replicas_per_master",
				Description: "The number of replicas to be created per master.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.ReplicasPerMaster"),
			},
			{
				Name:        "shard_count",
				Description: "The number of shards to be created on a premium cluster cache.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.ShardCount"),
			},
			{
				Name:        "access_keys",
				Description: "The keys of the Redis cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AccessKeys"),
			},
			{
				Name:        "linked_servers",
				Description: "A list of the linked servers associated with the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.LinkedServers"),
			},
			{
				Name:        "instances",
				Description: "A list of the Redis instances associated with the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Instances"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "A list of private endpoint connection associated with the specified redis cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PrivateEndpointConnections"),
			},
			{
				Name:        "redis_configuration",
				Description: "Describes the redis cache configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.RedisConfiguration"),
			},
			{
				Name:        "tenant_settings",
				Description: "A dictionary of tenant settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.TenantSettings"),
			},
			{
				Name:        "zones",
				Description: "A list of availability zones denoting where the resource needs to come from.",
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
				Name:        "environment_name",
				Description: ColumnDescriptionEnvironmentName,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getEnvironmentName).WithCache(),
				Transform:   transform.FromValue(),
			},
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

func listRedisCaches(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRedisCaches")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := redis.NewClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.ListBySubscription(ctx)
	if err != nil {
		return nil, err
	}

	for _, cache := range result.Values() {
		d.StreamListItem(ctx, cache)
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

		for _, cache := range result.Values() {
			d.StreamListItem(ctx, cache)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedisCache(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRedisCache")

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
	client := redis.NewClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
