package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type subnetInfo = struct {
	Subnet         network.Subnet
	Name           *string
	VirtualNetwork *string
	ResourceGroup  *string
}

//// TABLE DEFINITION ////

func tableAzureSubnet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_subnet",
		Description: "Azure Subnet",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "virtual_network_name", "resource_group"}),
			ItemFromKey:       subnetDataFromKey,
			Hydrate:           getSubnet,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "NotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listVirtualNetworks,
			Hydrate:       listSubnets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the subnet",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a subnet uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.ID"),
			},
			{
				Name:        "virtual_network_name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name of the virtual network in which the subnet is created",
				Transform:   transform.FromField("VirtualNetwork"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.Etag"),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
				Default:     "Microsoft.Network/subnets",
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the subnet resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "address_prefix",
				Description: "Contains the address prefix for the subnet",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.AddressPrefix"),
			},
			{
				Name:        "nat_gateway_id",
				Description: "The ID of the Nat gateway associated with the subnet",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.NatGateway.ID"),
			},
			{
				Name:        "network_security_group_id",
				Description: "Network security group associated with the subnet",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.NetworkSecurityGroup.ID"),
			},
			{
				Name:        "private_endpoint_network_policies",
				Description: "Enable or Disable apply network policies on private end point in the subnet",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.PrivateEndpointNetworkPolicies"),
			},
			{
				Name:        "private_link_service_network_policies",
				Description: "Enable or Disable apply network policies on private link service in the subnet",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.PrivateLinkServiceNetworkPolicies"),
			},
			{
				Name:        "route_table_id",
				Description: "Route table associated with the subnet",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.RouteTable.ID"),
			},
			{
				Name:        "delegations",
				Description: "A list of references to the delegations on the subnet",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.Delegations"),
			},
			{
				Name:        "service_endpoints",
				Description: "A list of service endpoints",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.ServiceEndpoints"),
			},
			{
				Name:        "service_endpoint_policies",
				Description: "A list of service endpoint policies",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.ServiceEndpointPolicies"),
			},

			// Standard columns
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
				Transform:   transform.FromField("Subnet.ID").Transform(idToAkas),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceGroup").Transform(toLower),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// BUILD HYDRATE INPUT ////

func subnetDataFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	resourceGroup := quals["resource_group"].GetStringValue()
	virtualNetwork := quals["virtual_network_name"].GetStringValue()
	item := &subnetInfo{
		Name:           &name,
		VirtualNetwork: &virtualNetwork,
		ResourceGroup:  &resourceGroup,
	}
	return item, nil
}

//// FETCH FUNCTIONS ////

func listSubnets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of virtual network
	virtualNetwork := h.Item.(network.VirtualNetwork)
	resourceGroupName := &strings.Split(string(*virtualNetwork.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	subnetClient := network.NewSubnetsClient(subscriptionID)
	subnetClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := subnetClient.List(ctx, *resourceGroupName, *virtualNetwork.Name)
		if err != nil {
			return nil, err
		}

		for _, subnet := range result.Values() {
			d.StreamLeafListItem(ctx, subnetInfo{subnet, subnet.Name, virtualNetwork.Name, resourceGroupName})
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS ////

func getSubnet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// defaultRegion := GetDefaultRegion()
	subnet := h.Item.(*subnetInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	subnetClient := network.NewSubnetsClient(subscriptionID)
	subnetClient.Authorizer = session.Authorizer

	op, err := subnetClient.Get(ctx, *subnet.ResourceGroup, *subnet.VirtualNetwork, *subnet.Name, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return subnetInfo{op, op.Name, subnet.VirtualNetwork, subnet.ResourceGroup}, nil
	}

	return nil, nil
}
