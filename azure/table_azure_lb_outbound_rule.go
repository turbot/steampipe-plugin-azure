package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureLoadBalancerOutboundRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lb_outbound_rule",
		Description: "Azure Load Balancer Outbound Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"load_balancer_name", "name", "resource_group"}),
			Hydrate:           getLoadBalancerOutboundRule,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate:       listLoadBalancerOutboundRules,
			ParentHydrate: listLoadBalancers,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource that is unique within the set of outbound rules used by the load balancer. This name can be used to access the resource.",
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
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the outbound rule resource. Possible values include: 'ProvisioningStateSucceeded', 'ProvisioningStateUpdating', 'ProvisioningStateDeleting', 'ProvisioningStateFailed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OutboundRulePropertiesFormat.ProvisioningState"),
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
				Name:        "allocated_outbound_ports",
				Description: "The number of outbound ports to be used for NAT.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("OutboundRulePropertiesFormat.AllocatedOutboundPorts"),
			},
			{
				Name:        "enable_tcp_reset",
				Description: "Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("OutboundRulePropertiesFormat.EnableTCPReset"),
				Default:     false,
			},
			{
				Name:        "idle_timeout_in_minutes",
				Description: "The timeout for the TCP idle connection. The value can be set between 4 and 30 minutes. The default value is 4 minutes. This element is only used when the protocol is set to TCP.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("OutboundRulePropertiesFormat.IdleTimeoutInMinutes"),
			},
			{
				Name:        "protocol",
				Description: "The protocol for the outbound rule in load balancer. Possible values include: 'LoadBalancerOutboundRuleProtocolTCP', 'LoadBalancerOutboundRuleProtocolUDP', 'LoadBalancerOutboundRuleProtocolAll'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OutboundRulePropertiesFormat.Protocol"),
			},
			{
				Name:        "backend_address_pools",
				Description: "A reference to a pool of DIPs. Outbound traffic is randomly load balanced across IPs in the backend IPs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OutboundRulePropertiesFormat.BackendAddressPools"),
			},
			{
				Name:        "frontend_ip_configurations",
				Description: "The Frontend IP addresses of the load balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OutboundRulePropertiesFormat.FrontendIPConfigurations"),
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
				Name:        "environment_name",
				Description: ColumnDescriptionEnvironmentName,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getEnvironmentName).WithCache(),
				Transform:   transform.FromValue(),
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

type LoadBalancerOutboundRulesInfo = struct {
	network.OutboundRule
	LoadBalancerName string
}

//// LIST FUNCTION

func listLoadBalancerOutboundRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of load balancer
	loadBalancer := h.Item.(network.LoadBalancer)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*loadBalancer.ID, "/")[4]

	listLoadBalancerOutboundClient := network.NewLoadBalancerOutboundRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	listLoadBalancerOutboundClient.Authorizer = session.Authorizer

	result, err := listLoadBalancerOutboundClient.List(ctx, resourceGroup, *loadBalancer.Name)
	if err != nil {
		return nil, err
	}
	for _, rule := range result.Values() {
		d.StreamListItem(ctx, LoadBalancerOutboundRulesInfo{rule, *loadBalancer.Name})
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
		for _, rule := range result.Values() {
			d.StreamListItem(ctx, LoadBalancerOutboundRulesInfo{rule, *loadBalancer.Name})
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

func getLoadBalancerOutboundRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoadBalancerOutboundRule")

	loadBalancerName := d.KeyColumnQuals["load_balancer_name"].GetStringValue()
	loadBalancerOutboundRuleName := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty check
	if loadBalancerName == "" || loadBalancerOutboundRuleName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	LoadBalancerOutboundRuleClient := network.NewLoadBalancerOutboundRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	LoadBalancerOutboundRuleClient.Authorizer = session.Authorizer

	op, err := LoadBalancerOutboundRuleClient.Get(ctx, resourceGroup, loadBalancerName, loadBalancerOutboundRuleName)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return LoadBalancerOutboundRulesInfo{op, loadBalancerName}, nil
	}

	return nil, nil
}
