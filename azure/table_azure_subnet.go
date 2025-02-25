package azure

import (
	"context"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
			KeyColumns: plugin.AllColumns([]string{"name", "virtual_network_name", "resource_group"}),
			Hydrate:    getSubnet,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "NotFound", "ResourceGroupNotFound", "404"}),
			},
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

	// Apply Retry rule
	ApplyRetryRules(ctx, &subnetClient, d.Connection)

	result, err := subnetClient.List(ctx, *resourceGroupName, *virtualNetwork.Name)
	if err != nil {
		return nil, err
	}
	for _, subnet := range result.Values() {
		d.StreamListItem(ctx, subnetInfo{subnet, subnet.Name, virtualNetwork.Name, resourceGroupName})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
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
			if d.RowsRemaining(ctx) == 0 {
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

	// Apply Retry rule
	ApplyRetryRules(ctx, &subnetClient, d.Connection)

	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	virtualNetwork := d.EqualsQuals["virtual_network_name"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()

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
	subnetClient := resources.NewClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	subnetClient.Authorizer = session.Authorizer

	var wg sync.WaitGroup
	ipCh := make(chan *map[string]interface{}, len(configurations))
	errorCh := make(chan error, len(configurations))

	for i := 0; i < len(configurations); i++ {
		wg.Add(1)
		go getIpConfigurationAsync(ctx, &configurations[i], subnetClient, &wg, ipCh, errorCh)
	}

	// wait for all ip configurations to be processed
	wg.Wait()

	// Close channels after all goroutines have finished
	close(ipCh)
	close(errorCh)

	// Collect any errors from errorCh
	var collectedErrors []error
	for err := range errorCh {
		collectedErrors = append(collectedErrors, err)
	}

	if len(collectedErrors) > 0 {
		// Return the first error encountered (or handle errors as needed)
		return nil, collectedErrors[0]
	}

	var ipConfigurations []*map[string]interface{}

	for ipConfig := range ipCh {
		ipConfigurations = append(ipConfigurations, ipConfig)
	}

	return ipConfigurations, nil
}

func getIpConfigurationAsync(ctx context.Context, ipConfig *network.IPConfiguration, client resources.Client, wg *sync.WaitGroup, ipCh chan *map[string]interface{}, errorCh chan error) {
	defer wg.Done()

	rowData, err := getIpConfiguration(ctx, ipConfig, client)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		ipCh <- rowData
	}
}

func getIpConfiguration(ctx context.Context, ipConfig *network.IPConfiguration, client resources.Client) (*map[string]interface{}, error) {
	if ipConfig == nil {
		return nil, nil
	}

	configurationId := *ipConfig.ID

	// API version is required to make the API call.
	// https://learn.microsoft.com/en-us/rest/api/resources/resources/get-by-id?view=rest-resources-2021-04-01#code-try-0
	apiVersion := "2021-04-01"

	configuration, err := client.GetByID(ctx, configurationId, apiVersion)
	if err != nil {
		plugin.Logger(ctx).Error("azure_subnet.getIpConfiguration", "api_error", err)
		return nil, err
	}

	resourceData := make(map[string]interface{})

	// Extract the properties unless the top-level properties are not being retrieved.
	if configuration.ID != nil {
		resourceData["id"] = *configuration.ID
	}
	if configuration.Plan != nil {
		resourceData["plan"] = *configuration.Plan
	}
	if configuration.Properties != nil {
		resourceData["properties"] = configuration.Properties
	}
	if configuration.Kind != nil {
		resourceData["kind"] = *configuration.Kind
	}
	if configuration.ManagedBy != nil {
		resourceData["managedBy"] = *configuration.ManagedBy
	}
	if configuration.Sku != nil {
		resourceData["sku"] = *configuration.Sku
	}
	if configuration.Identity != nil {
		resourceData["identity"] = *configuration.Identity
	}
	if configuration.Name != nil {
		resourceData["name"] = *configuration.Name
	}
	if configuration.Type != nil {
		resourceData["type"] = *configuration.Type
	}
	if configuration.Location != nil {
		resourceData["location"] = *configuration.Location
	}
	if configuration.ExtendedLocation != nil {
		resourceData["extendedLocation"] = *configuration.ExtendedLocation
	}
	if configuration.Tags != nil {
		resourceData["tags"] = configuration.Tags
	}

	return &resourceData, nil
}
