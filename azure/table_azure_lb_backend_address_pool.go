package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAzureLoadBalancerBackendAddressPool(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lb_backend_address_pool",
		Description: "Azure Load Balancer Backend Address Pool",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"load_balancer_name", "name", "resource_group"}),
			Hydrate:    getBackendAddressPool,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate:       listBackendAddressPools,
			ParentHydrate: listLoadBalancers,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource that is unique within the set of backend address pools used by the load balancer. This name can be used to access the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "load_balancer_name",
				Description: "The friendly name that identifies the load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractLoadBalancerNameFromBackendAddressPoolID),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the backend address pool resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackendAddressPoolPropertiesFormat.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "outbound_rule_id",
				Description: "A reference to an outbound rule that uses this backend address pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackendAddressPoolPropertiesFormat.OutboundRule.ID"),
			},
			{
				Name:        "backend_ip_configurations",
				Description: "An array of references to IP addresses defined in network interfaces.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendAddressPoolPropertiesFormat.BackendIPConfigurations"),
			},
			{
				Name:        "gateway_load_balancer_tunnel_interface",
				Description: "An array of gateway load balancer tunnel interfaces.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendAddressPoolPropertiesFormat.GatewayLoadBalancerTunnelInterface"),
			},
			{
				Name:        "load_balancer_backend_addresses",
				Description: "An array of backend addresses.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses"),
			},
			{
				Name:        "load_balancing_rules",
				Description: "An array of references to load balancing rules that use this backend address pool.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendAddressPoolPropertiesFormat.LoadBalancingRules"),
			},
			{
				Name:        "outbound_rules",
				Description: "An array of references to outbound rules that use this backend address pool.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendAddressPoolPropertiesFormat.OutboundRules"),
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

			// Azure standard columns
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listBackendAddressPools(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of load balancer
	loadBalancer := h.Item.(network.LoadBalancer)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*loadBalancer.ID, "/")[4]

	listBackendAddressPoolsClient := network.NewLoadBalancerBackendAddressPoolsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	listBackendAddressPoolsClient.Authorizer = session.Authorizer

	result, err := listBackendAddressPoolsClient.List(ctx, resourceGroup, *loadBalancer.Name)
	if err != nil {
		return nil, err
	}
	for _, backendAddressPool := range result.Values() {
		d.StreamListItem(ctx, backendAddressPool)
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
		for _, backendAddressPool := range result.Values() {
			d.StreamListItem(ctx, backendAddressPool)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getBackendAddressPool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBackendAddressPool")

	loadBalancerName := d.KeyColumnQuals["load_balancer_name"].GetStringValue()
	backendAddressPoolName := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty loadBalancerName, backendAddressPoolName or resourceGroup
	if loadBalancerName == "" || backendAddressPoolName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	BackendAddressPoolClient := network.NewLoadBalancerBackendAddressPoolsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	BackendAddressPoolClient.Authorizer = session.Authorizer

	op, err := BackendAddressPoolClient.Get(ctx, resourceGroup, loadBalancerName, backendAddressPoolName)
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

//// TRANSFORM FUNCTION

func extractLoadBalancerNameFromBackendAddressPoolID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(network.BackendAddressPool)
	loadBalancerName := strings.Split(*data.ID, "/")[8]
	return loadBalancerName, nil
}
