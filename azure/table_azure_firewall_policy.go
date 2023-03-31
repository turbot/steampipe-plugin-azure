package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureFirewallPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_firewall_policy",
		Description: "Azure Firewall Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getFirewallPolicy,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listFirewallPolicies,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the firewall policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a firewall policy uniquely.",
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
				Description: "The resource type of the firewall policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the firewall policy resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FirewallPolicyPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "intrusion_detection_mode",
				Description: "Intrusion detection general state. Possible values include: 'FirewallPolicyIntrusionDetectionStateTypeOff', 'FirewallPolicyIntrusionDetectionStateTypeAlert', 'FirewallPolicyIntrusionDetectionStateTypeDeny'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.IntrusionDetection.Mode"),
			},
						{
				Name:        "sku_tier",
				Description: "Tier of Firewall Policy. Possible values include: 'FirewallPolicySkuTierStandard', 'FirewallPolicySkuTierPremium'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.Sku.Tier"),
			},
			{
				Name:        "threat_intel_mode",
				Description: "The operation mode for Threat Intelligence. Possible values include: 'AzureFirewallThreatIntelModeAlert', 'AzureFirewallThreatIntelModeDeny', 'AzureFirewallThreatIntelModeOff'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FirewallPolicyPropertiesFormat.ThreatIntelMode.ID"),
			},
			{
				Name:        "base_policy",
				Description: "The parent firewall policy from which rules are inherited.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.BasePolicy"),
			},
			{
				Name:        "child_policies",
				Description: "List of references to Child Firewall Policies.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.ChildPolicies"),
			},
			{
				Name:        "dns_settings",
				Description: "DNS Proxy Settings definition.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.DNSSettings"),
			},
			{
				Name:        "firewalls",
				Description: "List of references to Azure Firewalls that this Firewall Policy is associated with.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.Firewalls"),
			},
			{
				Name:        "identity",
				Description: "The identity of the firewall policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "intrusion_detection_configuration",
				Description: "Intrusion detection configuration properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.IntrusionDetection.Configuration"),
			},
			{
				Name:        "rule_collection_groups",
				Description: "List of references to FirewallPolicyRuleCollectionGroups.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.RuleCollectionGroups"),
			},
			{
				Name:        "threat_intel_whitelist_ip_addresses",
				Description: "List of IP addresses for the ThreatIntel Whitelist.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.IPAddresses"),
			},
			{
				Name:        "threat_intel_whitelist_fqdns",
				Description: "List of FQDNs for the ThreatIntel Whitelist.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.Fqdns"),
			},
			{
				Name:        "transport_security_certificate_authority",
				Description: "The CA used for intermediate CA generation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.TransportSecurity.CertificateAuthority"),
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

func listFirewallPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_firewall_policy.listFirewallPolicies", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewFirewallPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer
	result, err := networkClient.ListAll(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_firewall_policy.listFirewallPolicies", "api_error", err)
		return nil, err
	}

	for _, policy := range result.Values() {
		d.StreamListItem(ctx, policy)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err := result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, policy := range result.Values() {
			d.StreamListItem(ctx, policy)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS ////

func getFirewallPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_firewall_policy.getFirewallPolicy", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewFirewallPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	op, err := networkClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("azure_firewall_policy.getFirewallPolicy", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
