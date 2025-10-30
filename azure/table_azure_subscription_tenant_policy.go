package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureSubscriptionTenantPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_subscription_tenant_policy",
		Description: "Azure Subscription Tenant Policy",
		List: &plugin.ListConfig{
			Hydrate: listSubscriptionTenantPolicy,
			Tags: map[string]string{
				"service": "Microsoft.Subscription",
				"action":  "policies/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Policy name.",
			},
			{
				Name:        "id",
				Description: "Policy ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_id",
				Description: "The policy ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PolicyID"),
			},
			{
				Name:        "block_subscriptions_leaving_tenant",
				Description: "Blocks the leaving of subscriptions from user's tenant.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.BlockSubscriptionsLeavingTenant"),
			},
			{
				Name:        "block_subscriptions_into_tenant",
				Description: "Blocks the entering of subscriptions into user's tenant.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.BlockSubscriptionsIntoTenant"),
			},
			{
				Name:        "exempted_principals",
				Description: "List of user objectIds that are exempted from the set subscription tenant policies for the user's tenant.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ExemptedPrincipals"),
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

func listSubscriptionTenantPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get the session with credentials
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "session_error", err)
		return nil, err
	}

	// Create the policy client
	client, err := armsubscription.NewPolicyClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "client_error", err)
		return nil, err
	}

	// Get the tenant policy
	result, err := client.GetPolicyForTenant(ctx, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "api_error", err)
		return nil, err
	}

	// Stream the result
	d.StreamListItem(ctx, result.GetTenantPolicyResponse)

	return nil, nil
}
