package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureFirewall(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_firewall",
		Description: "Azure Firewall",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getFirewall,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listFirewalls,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the firewall",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a firewall uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the firewall",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the firewall resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "firewall_policy_id",
				Description: "The firewallPolicy associated with this azure firewall",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.FirewallPolicy.ID"),
			},
			{
				Name:        "hub_private_ip_address",
				Description: "Private IP Address associated with azure firewall",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.HubIPAddresses.PrivateIPAddress"),
			},
			{
				Name:        "hub_public_ip_address_count",
				Description: "The number of Public IP addresses associated with azure firewall",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.HubIPAddresses.PublicIPs.Count"),
			},
			{
				Name:        "sku_name",
				Description: "Name of an Azure Firewall SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Tier of an Azure Firewall",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "threat_intel_mode",
				Description: "The operation mode for Threat Intelligence",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.ThreatIntelMode").Transform(transform.ToString),
			},
			{
				Name:        "virtual_hub_id",
				Description: "The virtualHub to which the firewall belongs",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.VirtualHub.ID"),
			},
			{
				Name:        "additional_properties",
				Description: "A collection of additional properties used to further config this azure firewall",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.AdditionalProperties"),
			},
			{
				Name:        "application_rule_collections",
				Description: "A collection of application rule collections used by Azure Firewall",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.ApplicationRuleCollections"),
			},
			{
				Name:        "availability_zones",
				Description: "A collection of availability zones denoting where the resource needs to come from",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Zones"),
			},
			{
				Name:        "hub_public_ip_addresses",
				Description: "A collection of Public IP addresses associated with azure firewall or IP addresses to be retained",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.HubIPAddresses.PublicIPs.Addresses"),
			},
			{
				Name:        "ip_configurations",
				Description: "A collection of IP configuration of the Azure Firewall resource",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(ipConfigurationData),
			},
			{
				Name:        "ip_groups",
				Description: "A collection of IpGroups associated with AzureFirewall",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.IPGroups"),
			},
			{
				Name:        "nat_rule_collections",
				Description: "A collection of NAT rule collections used by Azure Firewall",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.NatRuleCollections"),
			},
			{
				Name:        "network_rule_collections",
				Description: "A collection of network rule collections used by Azure Firewall",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFirewallPropertiesFormat.NetworkRuleCollections"),
			},

			// Standard columns
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// FETCH FUNCTIONS ////

func listFirewalls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewAzureFirewallsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer
	result, err := networkClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, firewall := range result.Values() {
		d.StreamListItem(ctx, firewall)
	}

	for result.NotDone() {
		err := result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, firewall := range result.Values() {
			d.StreamListItem(ctx, firewall)
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getFirewall(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFirewall")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewAzureFirewallsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	op, err := networkClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

//// Transform Functions

func ipConfigurationData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(network.AzureFirewall)

	var output []map[string]interface{}
	for _, firewall := range *data.AzureFirewallPropertiesFormat.IPConfigurations {
		objectMap := make(map[string]interface{})
		if firewall.AzureFirewallIPConfigurationPropertiesFormat.PrivateIPAddress != nil {
			objectMap["privateIPAddress"] = firewall.AzureFirewallIPConfigurationPropertiesFormat.PrivateIPAddress
		}
		if firewall.AzureFirewallIPConfigurationPropertiesFormat.PublicIPAddress != nil {
			objectMap["publicIPAddress"] = firewall.AzureFirewallIPConfigurationPropertiesFormat.PublicIPAddress
		}
		if firewall.AzureFirewallIPConfigurationPropertiesFormat.Subnet != nil {
			objectMap["subnet"] = firewall.AzureFirewallIPConfigurationPropertiesFormat.Subnet
		}
		if firewall.AzureFirewallIPConfigurationPropertiesFormat.ProvisioningState != "" {
			objectMap["provisioningState"] = firewall.AzureFirewallIPConfigurationPropertiesFormat.ProvisioningState
		}
		output = append(output, objectMap)
	}
	return output, nil
}
