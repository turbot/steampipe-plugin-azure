package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureLoadBalancerRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lb_rule",
		Description: "Azure Load Balancer Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"load_balancer_name", "name", "resource_group"}),
			Hydrate:    getLoadBalancerRule,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate:       listLoadBalancerRules,
			ParentHydrate: listLoadBalancers,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource that is unique within the set of load balancing rules used by the load balancer. This name can be used to access the resource.",
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
				Transform:   transform.From(extractLoadBalancerNameFromRuleID),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the load balancing rule resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.ProvisioningState"),
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
				Name:        "backend_address_pool_id",
				Description: "A reference to a pool of DIPs. Inbound traffic is randomly load balanced across IPs in the backend IPs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.BackendAddressPool.ID"),
			},
			{
				Name:        "backend_port",
				Description: "The port used for internal connections on the endpoint. Acceptable values are between 0 and 65535. Note that value 0 enables 'Any Port'.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.BackendPort"),
			},
			{
				Name:        "disable_outbound_snat",
				Description: "Configures SNAT for the VMs in the backend pool to use the publicIP address specified in the frontend of the load balancing rule.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.DisableOutboundSnat"),
				Default:     false,
			},
			{
				Name:        "enable_floating_ip",
				Description: "Configures a virtual machine's endpoint for the floating IP capability required to configure a SQL AlwaysOn Availability Group. This setting is required when using the SQL AlwaysOn Availability Groups in SQL server. This setting can't be changed after you create the endpoint.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.EnableFloatingIP"),
				Default:     false,
			},
			{
				Name:        "enable_tcp_reset",
				Description: "Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.EnableTCPReset"),
				Default:     false,
			},
			{
				Name:        "frontend_ip_configuration_id",
				Description: "A reference to frontend IP addresses.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.FrontendIPConfiguration.ID"),
			},
			{
				Name:        "frontend_port",
				Description: "The port for the external endpoint. Port numbers for each rule must be unique within the Load Balancer. Acceptable values are between 0 and 65534. Note that value 0 enables 'Any Port'.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.FrontendPort"),
			},
			{
				Name:        "idle_timeout_in_minutes",
				Description: "The timeout for the TCP idle connection. The value can be set between 4 and 30 minutes. The default value is 4 minutes. This element is only used when the protocol is set to TCP.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.IdleTimeoutInMinutes"),
			},
			{
				Name:        "load_distribution",
				Description: "The load distribution policy for this rule. Possible values include: 'Default', 'SourceIP', 'SourceIPProtocol'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.LoadDistribution"),
			},
			{
				Name:        "probe_id",
				Description: "The reference to the load balancer probe used by the load balancing rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.Probe.ID"),
			},
			{
				Name:        "protocol",
				Description: "The reference to the transport protocol used by the load balancing rule. Possible values include: 'UDP', 'TCP', 'All'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.Protocol"),
			},
			{
				Name:        "backend_address_pools",
				Description: "An array of references to pool of DIPs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancingRulePropertiesFormat.BackendAddressPools"),
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

func listLoadBalancerRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of load balancer
	loadBalancer := h.Item.(network.LoadBalancer)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*loadBalancer.ID, "/")[4]

	listLoadBalancerRulesClient := network.NewLoadBalancerLoadBalancingRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	listLoadBalancerRulesClient.Authorizer = session.Authorizer

	result, err := listLoadBalancerRulesClient.List(ctx, resourceGroup, *loadBalancer.Name)
	if err != nil {
		return nil, err
	}
	for _, rule := range result.Values() {
		d.StreamListItem(ctx, rule)
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
		for _, rule := range result.Values() {
			d.StreamListItem(ctx, rule)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getLoadBalancerRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoadBalancerRule")

	loadBalancerName := d.EqualsQuals["load_balancer_name"].GetStringValue()
	loadBalancerRuleName := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Handle empty loadBalancerName, loadBalancerRuleName or resourceGroup
	if loadBalancerName == "" || loadBalancerRuleName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	LoadBalancerRuleClient := network.NewLoadBalancerLoadBalancingRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	LoadBalancerRuleClient.Authorizer = session.Authorizer

	op, err := LoadBalancerRuleClient.Get(ctx, resourceGroup, loadBalancerName, loadBalancerRuleName)
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

func extractLoadBalancerNameFromRuleID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(network.LoadBalancingRule)
	loadBalancerName := strings.Split(*data.ID, "/")[8]
	return loadBalancerName, nil
}
