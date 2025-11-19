package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureTenant(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_tenant",
		Description: "Azure Tenant",
		List: &plugin.ListConfig{
			Hydrate: listTenants,
			Tags: map[string]string{
				"service": "Microsoft.Resources",
				"action":  "tenants/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name of the tenant.",
				Transform:   transform.From(getNameOrID),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The fully qualified ID of the tenant. For example, /tenants/00000000-0000-0000-0000-000000000000.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "tenant_id",
				Type:        proto.ColumnType_STRING,
				Description: "The tenant ID. For example, 00000000-0000-0000-0000-000000000000.",
				Transform:   transform.FromField("TenantID"),
			},
			{
				Name:        "tenant_category",
				Type:        proto.ColumnType_STRING,
				Description: "The tenant category. Possible values include: 'Home', 'ProjectedBy', 'ManagedBy'.",
				Transform:   transform.FromField("TenantCategory").Transform(transform.ToString),
			},
			{
				Name:        "country",
				Type:        proto.ColumnType_STRING,
				Description: "Country/region name of the address for the tenant.",
			},
			{
				Name:        "country_code",
				Type:        proto.ColumnType_STRING,
				Description: "Country/region abbreviation for the tenant.",
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The list of domains for the tenant.",
			},
		{
			Name:        "domains",
			Type:        proto.ColumnType_JSON,
			Description: "The list of domains for the tenant.",
		},
		{
			Name:        "subscription_policy",
			Type:        proto.ColumnType_JSON,
			Description: "The subscription policy for the tenant, including properties like BlockSubscriptionsLeavingTenant, BlockSubscriptionsIntoTenant, and ExemptedPrincipals.",
			Hydrate:     getTenantSubscriptionPolicy,
			Transform:   transform.FromValue(),
		},

		// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getNameOrID),
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

func listTenants(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	client := subscriptions.NewTenantsClientWithBaseURI(session.ResourceManagerEndpoint)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, resp := range op.Values() {
		d.StreamListItem(ctx, resp)
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getTenantSubscriptionPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get the session with credentials
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_tenant.getTenantSubscriptionPolicy", "session_error", err)
		return nil, err
	}

	// Create the policy client
	client, err := armsubscription.NewPolicyClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_tenant.getTenantSubscriptionPolicy", "client_error", err)
		return nil, err
	}

	// Get the tenant policy
	result, err := client.GetPolicyForTenant(ctx, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_tenant.getTenantSubscriptionPolicy", "api_error", err)
		return nil, err
	}

	return result.GetTenantPolicyResponse, nil
}

//// TRANSFORM FUNCTION

func getNameOrID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(subscriptions.TenantIDDescription)
	if data.DisplayName != nil {
		return data.DisplayName, nil
	}
	return data.TenantID, nil
}
