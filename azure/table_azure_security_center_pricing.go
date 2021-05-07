package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterPricing(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_pricing",
		Description: "Azure Security Center Pricing",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getSecurityCenterPricing,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterPricings,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The resource id.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pricing_tier",
				Description: "Pricing tier type. Possible values include: 'Free', 'Standard'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PricingProperties.PricingTier"),
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

			// Azure standard columns
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

func listSecurityCenterPricings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	pricingClient := security.NewPricingsClient(subscriptionID, "")
	pricingClient.Authorizer = session.Authorizer

	pricingList, err := pricingClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, pricing := range pricingList.Values() {
		d.StreamListItem(ctx, pricing)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCenterPricing(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	pricingClient := security.NewPricingsClient(subscriptionID, "")
	pricingClient.Authorizer = session.Authorizer

	pricing, err := pricingClient.GetSubscriptionPricing(ctx, name)
	if err != nil {
		return err, nil
	}

	return pricing, nil
}
