package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
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
			ParentHydrate: listResourceGroups,
			Hydrate:       listVirtualNetworkGateways,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the virtual network gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a virtual network gateway uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "gateway_type",
				Description: "The type of this virtual network gateway. Possible values include: 'Vpn', 'ExpressRoute'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.GatewayType").Transform(transform.ToString),
			},
			{
				Name:        "vpn_type",
				Description: "The type of this virtual network gateway. Valid values are: 'PolicyBased', 'RouteBased'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.VpnType").Transform(transform.ToString),
			},
			{
				Name:        "vpn_gateway_generation",
				Description: "The generation for this virtual network gateway. Must be None if gatewayType is not VPN. Valid values are: 'None', 'Generation1', 'Generation2'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.VpnGatewayGeneration").Transform(transform.ToString),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the virtual network gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "active_active",
				Description: "Indicates whether virtual network gateway configured with active-active mode, or not. If true, each Azure gateway instance will have a unique public IP address, and each will establish an IPsec/IKE S2S VPN tunnel to your on-premises VPN device specified in your local network gateway and connection.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.ActiveActive"),
			},
			{
				Name:        "enable_bgp",
				Description: "Indicates whether BGP is enabled for this virtual network gateway, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.EnableBgp"),
			},
			{
				Name:        "enable_dns_forwarding",
				Description: "Indicates whether DNS forwarding is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.EnableVMProtection"),
			},
			{
				Name:        "enable_private_ip_address",
				Description: "Indicates whether private IP needs to be enabled on this gateway for connections or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.EnablePrivateIPAddress"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "gateway_default_site",
				Description: "The reference to the LocalNetworkGateway resource, which represents local network site having default routes. Assign Null value in case of removing existing default site setting.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.GatewayDefaultSite.ID"),
			},
			{
				Name:        "inbound_dns_forwarding_endpoint",
				Description: "The IP address allocated by the gateway to which dns requests can be sent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.InboundDNSForwardingEndpoint"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the virtual network gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "sku_name",
				Description: "Gateway SKU name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Gateway SKU tier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "sku_capacity",
				Description: "Gateway SKU capacity.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.Sku.Capacity"),
			},
			{
				Name:        "bgp_settings",
				Description: "Virtual network gateway's BGP speaker settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.BgpSettings"),
			},
			{
				Name:        "custom_routes_address_prefixes",
				Description: "A list of address blocks reserved for this virtual network in CIDR notation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkGatewayPropertiesFormat.CustomRoutes.AddressPrefixes"),
			},
			{
				Name:        "gateway_connections",
				Description: "A list of virtual network gateway connection resources that exists in a resource group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVirtualNetworkGatewayConnection,
				Transform:   transform.FromValue(),
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

	data := h.Item.(resources.Group)
	resourceGroupName := *data.Name

	pagesLeft := true
	for pagesLeft {
		result, err := networkClient.List(ctx, resourceGroupName)
		if err != nil {
			return nil, err
		}

		for _, networkGateway := range result.Values() {
			d.StreamListItem(ctx, networkGateway)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVirtualNetworkGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVirtualNetworkGateway")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	networkClient := network.NewVirtualNetworkGatewaysClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	op, err := networkClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getVirtualNetworkGatewayConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVirtualNetworkGatewayConnection")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	virtualNetworkGateway := h.Item.(network.VirtualNetworkGateway)
	name := *virtualNetworkGateway.Name
	resourceGroup := strings.Split(*virtualNetworkGateway.ID, "/")[4]

	networkClient := network.NewVirtualNetworkGatewaysClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	var gatewayConnections []network.VirtualNetworkGatewayConnectionListEntity
	pagesLeft := true
	for pagesLeft {
		result, err := networkClient.ListConnections(ctx, resourceGroup, name)
		if err != nil {
			return nil, err
		}
		gatewayConnections = append(gatewayConnections, result.Values()...)
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return gatewayConnections, nil
}
