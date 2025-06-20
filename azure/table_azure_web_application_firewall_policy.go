package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureWebApplicationFirewallPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_web_application_firewall_policy",
		Description: "Azure Web Application Firewall Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getWebApplicationFirewallPolicy,
			Tags: map[string]string{
				"service": "Microsoft.Network",
				"action":  "webApplicationFirewallPolicies/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listWebApplicationFirewallPolicies,
			Tags: map[string]string{
				"service": "Microsoft.Network",
				"action":  "webApplicationFirewallPolicies/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Resource name.",
			},
			{
				Name:        "id",
				Description: "Resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the web application firewall policy resource. Possible values include: 'ProvisioningStateSucceeded', 'ProvisioningStateUpdating', 'ProvisioningStateDeleting', 'ProvisioningStateFailed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.ProvisioningState"),
			},
			{
				Name:        "resource_state",
				Description: "Resource status of the policy. Possible values include: 'WebApplicationFirewallPolicyResourceStateCreating', 'WebApplicationFirewallPolicyResourceStateEnabling', 'WebApplicationFirewallPolicyResourceStateEnabled', 'WebApplicationFirewallPolicyResourceStateDisabling', 'WebApplicationFirewallPolicyResourceStateDisabled', 'WebApplicationFirewallPolicyResourceStateDeleting'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.ResourceState"),
			},
			{
				Name:        "policy_settings",
				Description: "The PolicySettings for policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.PolicySettings"),
			},
			{
				Name:        "custom_rules",
				Description: "The custom rules inside the policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.CustomRules"),
			},
			{
				Name:        "application_gateways",
				Description: "A collection of references to application gateways.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.ApplicationGateways"),
			},
			{
				Name:        "managed_rules",
				Description: "Describes the managedRules structure.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.ManagedRules"),
			},
			{
				Name:        "http_listeners",
				Description: "A collection of references to application gateway http listeners.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.HTTPListeners"),
			},
			{
				Name:        "path_based_rules",
				Description: "A collection of references to application gateway path rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WebApplicationFirewallPolicyPropertiesFormat.PathBasedRules"),
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listWebApplicationFirewallPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_web_application_firewall_policy.listWebApplicationFirewallPolicies", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	networkClient := network.NewWebApplicationFirewallPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkClient, d.Connection)

	result, err := networkClient.ListAll(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_web_application_firewall_policy.listWebApplicationFirewallPolicies", "api_error", err)
		return nil, err
	}
	for _, firewallPolicy := range result.Values() {
		d.StreamListItem(ctx, firewallPolicy)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_web_application_firewall_policy.listWebApplicationFirewallPolicies", "paging_error", err)
			return nil, err
		}

		for _, firewallPolicy := range result.Values() {
			d.StreamListItem(ctx, firewallPolicy)
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

func getWebApplicationFirewallPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_web_application_firewall_policy.getWebApplicationFirewallPolicy", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewWebApplicationFirewallPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkClient, d.Connection)

	op, err := networkClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_web_application_firewall_policy.getWebApplicationFirewallPolicy", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
