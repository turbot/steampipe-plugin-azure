package azure

import (
	"context"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

type subnetInfo = struct {
	Subnet         network.Subnet
	Name           *string
	VirtualNetwork *string
	ResourceGroup  *string
}

//// TABLE DEFINITION

func tableAzureSubnet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_subnet",
		Description: "Azure Subnet",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "virtual_network_name", "resource_group"}),
			Hydrate:           getSubnet,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "NotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listVirtualNetworks,
			Hydrate:       listSubnets,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the subnet.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a subnet uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.ID"),
			},
			{
				Name:        "virtual_network_name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name of the virtual network in which the subnet is created.",
				Transform:   transform.FromField("VirtualNetwork"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.Etag"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Default:     "Microsoft.Network/subnets",
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the subnet resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "address_prefix",
				Description: "Contains the address prefix for the subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.AddressPrefix"),
			},
			{
				Name:        "nat_gateway_id",
				Description: "The ID of the Nat gateway associated with the subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.NatGateway.ID"),
			},
			{
				Name:        "network_security_group_id",
				Description: "Network security group associated with the subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.NetworkSecurityGroup.ID"),
			},
			{
				Name:        "private_endpoint_network_policies",
				Description: "Enable or Disable apply network policies on private end point in the subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.PrivateEndpointNetworkPolicies"),
			},
			{
				Name:        "private_link_service_network_policies",
				Description: "Enable or Disable apply network policies on private link service in the subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.PrivateLinkServiceNetworkPolicies"),
			},
			{
				Name:        "route_table_id",
				Description: "Route table associated with the subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.RouteTable.ID"),
			},
			{
				Name:        "delegations",
				Description: "A list of references to the delegations on the subnet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.Delegations"),
			},
			{
				Name:        "ip_configurations",
				Description: "IP Configuration details in a subnet.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSubnetIpConfigurations,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "service_endpoints",
				Description: "A list of service endpoints.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.ServiceEndpoints"),
			},
			{
				Name:        "service_endpoint_policies",
				Description: "A list of service endpoint policies.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Subnet.SubnetPropertiesFormat.ServiceEndpointPolicies"),
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
				Transform:   transform.FromField("Subnet.ID").Transform(idToAkas),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceGroup").Transform(toLower),
			},
		}),
	}
}

//// LIST FUNCTION

func listSubnets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of virtual network
	virtualNetwork := h.Item.(network.VirtualNetwork)
	resourceGroupName := &strings.Split(string(*virtualNetwork.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	subnetClient := network.NewSubnetsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	subnetClient.Authorizer = session.Authorizer

	result, err := subnetClient.List(ctx, *resourceGroupName, *virtualNetwork.Name)
	if err != nil {
		return nil, err
	}
	for _, subnet := range result.Values() {
		d.StreamListItem(ctx, subnetInfo{subnet, subnet.Name, virtualNetwork.Name, resourceGroupName})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, subnet := range result.Values() {
			d.StreamListItem(ctx, subnetInfo{subnet, subnet.Name, virtualNetwork.Name, resourceGroupName})
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSubnet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSubnet")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	subnetClient := network.NewSubnetsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	subnetClient.Authorizer = session.Authorizer

	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	virtualNetwork := d.KeyColumnQuals["virtual_network_name"].GetStringValue()
	name := d.KeyColumnQuals["name"].GetStringValue()

	op, err := subnetClient.Get(ctx, resourceGroup, virtualNetwork, name, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return subnetInfo{op, op.Name, &virtualNetwork, &resourceGroup}, nil
	}

	return nil, nil
}

// List or Get call of subnet doesn't return more info about ip configuration except the ip configuration id.
func getSubnetIpConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSubnetIpConfigurations")
	subnet := h.Item.(subnetInfo).Subnet

	configurations := []network.IPConfiguration{}

	if subnet.SubnetPropertiesFormat != nil {
		if subnet.SubnetPropertiesFormat.IPConfigurations != nil {
			configurations = *subnet.SubnetPropertiesFormat.IPConfigurations
		}
	}
	if len(configurations) <= 0 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	subnetClient := network.NewInterfaceIPConfigurationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	subnetClient.Authorizer = session.Authorizer

	var wg sync.WaitGroup
	ipCh := make(chan *network.InterfaceIPConfiguration, len(configurations))
	errorCh := make(chan error, len(configurations))

	for i := 0; i < len(configurations); i++ {
		wg.Add(1)
		go getIpConfigurationAsync(ctx, &configurations[i], subnetClient, &wg, ipCh, errorCh)
	}

	// wait for all ip configurations to be processed
	wg.Wait()
	// NOTE: close channel before ranging over results
	close(ipCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		return nil, err
	}

	var ipConfigurations []*network.InterfaceIPConfiguration

	for ipConfig := range ipCh {
		ipConfigurations = append(ipConfigurations, ipConfig)
	}

	return ipConfigurations, nil
}

func getIpConfigurationAsync(ctx context.Context, ipConfig *network.IPConfiguration, client network.InterfaceIPConfigurationsClient, wg *sync.WaitGroup, ipCh chan *network.InterfaceIPConfiguration, errorCh chan error) {
	defer wg.Done()

	rowData, err := getIpConfiguration(ctx, ipConfig, client)
	if err != nil {
		errorCh <- err
	} else if rowData.ID != nil {
		ipCh <- rowData
	}
}

func getIpConfiguration(ctx context.Context, ipConfig *network.IPConfiguration, client network.InterfaceIPConfigurationsClient) (*network.InterfaceIPConfiguration, error) {
	if ipConfig == nil {
		return nil, nil
	}

	configurationId := *ipConfig.ID
	resourceGroup := strings.Split(configurationId, "/")[4]
	networkInterface := strings.Split(configurationId, "/")[8]
	configName := strings.Split(configurationId, "/")[len(strings.Split(configurationId, "/"))-1]

	configuration, err := client.Get(ctx, resourceGroup, networkInterface, configName)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
