package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
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
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "session_error", err)
		return nil, err
	}

	// Build the request URL
	apiVersion := "2021-10-01"
	url := fmt.Sprintf("%s/providers/Microsoft.Subscription/policies/default?api-version=%s", session.ResourceManagerEndpoint, apiVersion)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "prepare_error", err)
		return nil, err
	}

	// Add authorization header
	preparer := autorest.CreatePreparer(session.Authorizer.WithAuthorization())
	req, err = preparer.Prepare(req)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "auth_error", err)
		return nil, err
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "read_error", err)
		return nil, err
	}

	// Parse the response
	var policy GetTenantPolicyResponse
	err = json.Unmarshal(body, &policy)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subscription_tenant_policy.listSubscriptionTenantPolicy", "unmarshal_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, policy)

	return nil, nil
}

//// HYDRATE FUNCTIONS

// Response structs for Subscription Tenant Policy API

// GetTenantPolicyResponse represents the tenant policy information
type GetTenantPolicyResponse struct {
	// ID - Policy Id
	ID *string `json:"id,omitempty"`
	// Name - Policy name
	Name *string `json:"name,omitempty"`
	// Type - Resource type
	Type *string `json:"type,omitempty"`
	// Properties - Tenant policy properties
	Properties *TenantPolicyProperties `json:"properties,omitempty"`
}

// TenantPolicyProperties represents tenant policy properties
type TenantPolicyProperties struct {
	// PolicyID - Policy Id
	PolicyID *string `json:"policyId,omitempty"`
	// BlockSubscriptionsLeavingTenant - Blocks the leaving of subscriptions from user's tenant
	BlockSubscriptionsLeavingTenant *bool `json:"blockSubscriptionsLeavingTenant,omitempty"`
	// BlockSubscriptionsIntoTenant - Blocks the entering of subscriptions into user's tenant
	BlockSubscriptionsIntoTenant *bool `json:"blockSubscriptionsIntoTenant,omitempty"`
	// ExemptedPrincipals - List of user objectIds that are exempted from the set subscription tenant policies for the user's tenant
	ExemptedPrincipals *[]string `json:"exemptedPrincipals,omitempty"`
}
