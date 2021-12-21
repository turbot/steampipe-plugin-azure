package azure

import (
	"context"

	sub "github.com/Azure/azure-sdk-for-go/profiles/latest/subscription/mgmt/subscription"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureLocation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_location",
		Description: "Azure Location",
		List: &plugin.ListConfig{
			Hydrate: listLocations,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The location name",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The display name of the location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The fully qualified ID of the location.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "latitude",
				Description: "The latitude of the location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "longitude",
				Description: "The longitude of the location",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
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
		},
	}
}

//// LIST FUNCTION

func listLocations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	subscriptionsClient := sub.NewSubscriptionsClientWithBaseURI(session.ResourceManagerEndpoint)
	subscriptionsClient.Authorizer = session.Authorizer

	result, err := subscriptionsClient.ListLocations(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	for _, location := range *result.Value {
		d.StreamListItem(ctx, location)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}
