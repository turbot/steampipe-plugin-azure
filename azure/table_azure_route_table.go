package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureRouteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_route_table",
		Description: "Azure Route Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getRouteTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listRouteTables,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the route table",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a route table uniquely",
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
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disable_bgp_route_propagation",
				Description: "Indicates Whether to disable the routes learned by BGP on that route table. True means disable.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("RouteTablePropertiesFormat.DisableBgpRoutePropagation"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the route table resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RouteTablePropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "routes",
				Description: "A list of routes contained within a route table",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RouteTablePropertiesFormat.Routes"),
			},
			{
				Name:        "subnets",
				Description: "A list of references to subnets",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RouteTablePropertiesFormat.Subnets"),
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID).Transform(toLower),
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listRouteTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	routeTableClient := network.NewRouteTablesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	routeTableClient.Authorizer = session.Authorizer

	result, err := routeTableClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, routeTable := range result.Values() {
		d.StreamListItem(ctx, routeTable)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, routeTable := range result.Values() {
			d.StreamListItem(ctx, routeTable)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getRouteTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRouteTable")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	routeTableClient := network.NewRouteTablesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	routeTableClient.Authorizer = session.Authorizer

	op, err := routeTableClient.Get(ctx, resourceGroup, name, "")
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
