package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_subscription",
		Description: "Azure Subscription",
		List: &plugin.ListConfig{
			Hydrate: listSubscriptions,
		},

		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The fully qualified ID for the subscription. For example, /subscriptions/00000000-0000-0000-0000-000000000000.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "subscription_id",
				Description: "The subscription ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubscriptionID"),
			},
			{
				Name:        "display_name",
				Description: "A friendly name that identifies a subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tenant_id",
				Description: "The subscription tenant ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TenantID"),
			},
			{
				Name:        "state",
				Description: "The subscription state. Possible values are Enabled, Warned, PastDue, Disabled, and Deleted. Possible values include: 'StateEnabled', 'StateWarned', 'StatePastDue', 'StateDisabled', 'StateDeleted'",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State").Transform(transform.ToString),
			},
			{
				Name:        "authorization_source",
				Description: "The authorization source of the request. Valid values are one or more combinations of Legacy, RoleBased, Bypassed, Direct and Management. For example, 'Legacy, RoleBased'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_by_tenants",
				Description: "An array containing the tenants managing the subscription.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subscription_policies",
				Description: "The subscription policies.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
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

func listSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	client := subscriptions.NewClient()
	client.Authorizer = session.Authorizer
	subscriptionID := session.SubscriptionID

	op, err := client.Get(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, op)

	return nil, nil
}
