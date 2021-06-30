package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureVirtualNetworkGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_virtual_network_gateway",
		Description: "Azure Virtual Network Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getVirtualNetworkGateway,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listVirtualNetworks,
			Hydrate:       listVirtualNetworkGateways,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the virtual network gateway.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a virtual network gateway uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_bgp",
				Description: "Whether BGP is enabled for this virtual network gateway or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.EnableBgp"),
			},
			{
				Name:        "enable_private_ip_address",
				Description: "Whether private IP needs to be enabled on this gateway for connections or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.EnablePrivateIPAddress"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the virtual network gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the virtual network gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "active_active",
				Description: "ActiveActive flag.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.ActiveActive"),
			},
			{
				Name:        "enable_dns_forwarding",
				Description: "Whether dns forwarding is enabled or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.EnableVMProtection"),
			},
			{
				Name:        "gateway_default_site",
				Description: "The reference to the LocalNetworkGateway resource which represents local network site having default routes. Assign Null value in case of removing existing default site setting.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.GatewayDefaultSite.ID"),
			},
			{
				Name:        "gateway_type",
				Description: "The type of this virtual network gateway. Possible values include: 'VirtualNetworkGatewayTypeVpn', 'VirtualNetworkGatewayTypeExpressRoute'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.GatewayType").Transform(transform.ToString),
			},
			{
				Name:        "inbound_dns_forwarding_endpoint",
				Description: "The IP address allocated by the gateway to which dns requests can be sent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.InboundDNSForwardingEndpoint"),
			},
			{
				Name:        "sku_name",
				Description: "Gateway SKU name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "Gateway SKU tier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.Sku.Tier"),
			},
			{
				Name:        "sku_capacity",
				Description: "Gateway SKU capacity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.Sku.Capacity"),
			},
			{
				Name:        "vpn_gateway_generation",
				Description: "The generation for this VirtualNetworkGateway. Must be None if gatewayType is not VPN. Possible values include: 'VpnGatewayGenerationNone', 'VpnGatewayGenerationGeneration1', 'VpnGatewayGenerationGeneration2'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.VpnGatewayGeneration").Transform(transform.ToString),
			},
			{
				Name:        "vpn_type",
				Description: "The type of this virtual network gateway. Possible values include: 'PolicyBased', 'RouteBased'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.VpnType"),
			},
			{
				Name:        "address_prefixes",
				Description: "A list of address blocks reserved for this virtual network in CIDR notation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.CustomRoutes.AddressPrefixes"),
			},
			{
				Name:        "bgp_settings",
				Description: "Virtual network gateway's BGP speaker settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.BgpSettings"),
			},
			{
				Name:        "ip_configurations",
				Description: "IP configurations for virtual network gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.IPConfigurations"),
			},
			{
				Name:        "vpn_client_configuration",
				Description: "The reference to the VpnClientConfiguration resource which represents the P2S VpnClient configurations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.VpnClientConfiguration"),
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

func listVirtualNetworkGateways(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewVirtualNetworkGatewaysClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	virtualNetwork := h.Item.(network.VirtualNetwork)
	resourceGroupName := strings.Split(string(*virtualNetwork.ID), "/")[4]

	pagesLeft := true
	for pagesLeft {
		result, err := networkClient.List(ctx, resourceGroupName)
		if err != nil {
			return nil, err
		}

		for _, network := range result.Values() {
			d.StreamListItem(ctx, network)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVirtualNetworkGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVirtualNetworkGateway")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewVirtualNetworkGatewaysClient(subscriptionID)
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
