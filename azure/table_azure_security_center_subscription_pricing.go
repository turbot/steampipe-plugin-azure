package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterPricing(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_subscription_pricing",
		Description: "Azure Security Center Subscription Pricing",
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
				Description: "The pricing id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "Name of the pricing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "pricing_tier",
				Type:        proto.ColumnType_STRING,
				Description: "The pricing tier value. Azure Security Center is provided in two pricing tiers: free and standard, with the standard tier available with a trial period. The standard tier offers advanced security capabilities, while the free tier offers basic security features.",
				Transform:   transform.FromField("PricingProperties.PricingTier"),
			},
			{
				Name:        "free_trial_remaining_time",
				Description: "The duration left for the subscriptions free trial period.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PricingProperties.FreeTrialRemainingTime"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of the pricing.",
				Transform:   transform.FromGo(),
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
	settingClient := security.NewPricingsClient(subscriptionID, "")
	settingClient.Authorizer = session.Authorizer

	pricingList, err := settingClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, pricing := range *pricingList.Value {
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

	// Handle empty input for get call
	if name == "" {
		return nil, nil
	}

	subscriptionID := session.SubscriptionID
	settingClient := security.NewPricingsClient(subscriptionID, name)
	settingClient.Authorizer = session.Authorizer

	setting, err := settingClient.Get(ctx, name)
	if err != nil {
		return err, nil
	}

	return setting, nil
}
