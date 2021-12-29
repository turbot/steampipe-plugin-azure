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

func tableAzureLoadBalancerNatRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lb_nat_rule",
		Description: "Azure Load Balancer Nat Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"load_balancer_name", "name", "resource_group"}),
			Hydrate:           getLoadBalancerNatRule,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate:       listLoadBalancerNatRules,
			ParentHydrate: listLoadBalancers,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource that is unique within the set of inbound NAT rules used by the load balancer. This name can be used to access the resource.",
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
				Description: "The provisioning state of the inbound NAT rule resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.ProvisioningState"),
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
				Name:        "backend_port",
				Description: "The port used for the internal endpoint. Acceptable values range from 1 to 65535.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.BackendPort"),
			},
			{
				Name:        "enable_floating_ip",
				Description: "Configures a virtual machine's endpoint for the floating IP capability required to configure a SQL AlwaysOn Availability Group. This setting is required when using the SQL AlwaysOn Availability Groups in SQL server. This setting can't be changed after you create the endpoint.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.EnableFloatingIP"),
				Default:     false,
			},
			{
				Name:        "frontend_port",
				Description: "The port for the external endpoint. Port numbers for each rule must be unique within the Load Balancer. Acceptable values range from 1 to 65534.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.FrontendPort"),
			},
			{
				Name:        "enable_tcp_reset",
				Description: "Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.EnableTCPReset"),
				Default:     false,
			},
			{
				Name:        "idle_timeout_in_minutes",
				Description: "The timeout for the TCP idle connection. The value can be set between 4 and 30 minutes. The default value is 4 minutes. This element is only used when the protocol is set to TCP.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.IdleTimeoutInMinutes"),
			},
			{
				Name:        "protocol",
				Description: "The reference to the transport protocol used by the load balancing rule. Possible values include: 'TransportProtocolUDP', 'TransportProtocolTCP', 'TransportProtocolAll'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.Protocol"),
			},
			{
				Name:        "backend_ip_configuration",
				Description: "A reference to a private IP address defined on a network interface of a VM. Traffic sent to the frontend port of each of the frontend IP configurations is forwarded to the backend IP.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.BackendIPConfiguration"),
			},
			{
				Name:        "frontend_ip_configuration",
				Description: "A reference to frontend IP addresses.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InboundNatRulePropertiesFormat.FrontendIPConfiguration"),
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

type LoadBalancerNatRulesInfo = struct {
	network.InboundNatRule
	LoadBalancerName string
}

//// LIST FUNCTION

func listLoadBalancerNatRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get load balancer details
	loadBalancer := h.Item.(network.LoadBalancer)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*loadBalancer.ID, "/")[4]

	natClient := network.NewInboundNatRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	natClient.Authorizer = session.Authorizer

	result, err := natClient.List(ctx, resourceGroup, *loadBalancer.Name)
	if err != nil {
		return nil, err
	}
	for _, rule := range result.Values() {
		d.StreamListItem(ctx, LoadBalancerNatRulesInfo{rule, *loadBalancer.Name})
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
			d.StreamListItem(ctx, LoadBalancerNatRulesInfo{rule, *loadBalancer.Name})
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

func getLoadBalancerNatRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoadBalancerNatRule")

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

	natClient := network.NewInboundNatRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	natClient.Authorizer = session.Authorizer

	op, err := natClient.Get(ctx, resourceGroup, loadBalancerName, loadBalancerOutboundRuleName, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return LoadBalancerNatRulesInfo{op, loadBalancerName}, nil
	}

	return nil, nil
}
