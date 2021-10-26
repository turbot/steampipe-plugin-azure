package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureNetworkWatcher(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_network_watcher",
		Description: "Azure Network Watcher",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getNetworkWatcher,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkWatchers,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the network watcher",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a network watcher uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the network watcher",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the network watcher resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WatcherPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},

			// Standard columns
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

//// FETCH FUNCTIONS ////

func listNetworkWatchers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkWatcherClient := network.NewWatchersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkWatcherClient.Authorizer = session.Authorizer

	result, err := networkWatcherClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, networkWatcher := range *result.Value {
		d.StreamListItem(ctx, networkWatcher)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getNetworkWatcher(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getNetworkWatcher")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkWatcherClient := network.NewWatchersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkWatcherClient.Authorizer = session.Authorizer

	op, err := networkWatcherClient.Get(ctx, resourceGroup, name)
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
