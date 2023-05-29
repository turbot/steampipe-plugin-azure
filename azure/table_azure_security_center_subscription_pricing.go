package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/security/mgmt/security"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
		Columns: azureColumns([]*plugin.Column{
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
		}),
	}
}

//// LIST FUNCTION

func listSecurityCenterPricings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	settingClient := security.NewPricingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	settingClient.Authorizer = session.Authorizer

	result, err := settingClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, pricing := range *result.Value {
		d.StreamListItem(ctx, pricing)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCenterPricing(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.EqualsQuals["name"].GetStringValue()

	// Handle empty input for get call
	if name == "" {
		return nil, nil
	}

	subscriptionID := session.SubscriptionID
	settingClient := security.NewPricingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	settingClient.Authorizer = session.Authorizer

	setting, err := settingClient.Get(ctx, name)
	if err != nil {
		return err, nil
	}

	return setting, nil
}
