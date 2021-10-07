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

func tableAzureLoadBalancerProbe(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lb_probe",
		Description: "Azure Load Balancer Probe",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"load_balancer_name", "name", "resource_group"}),
			Hydrate:           getLoadBalancerProbe,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate:       listLoadBalancerProbes,
			ParentHydrate: listLoadBalancers,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource that is unique within the set of probes used by the load balancer. This name can be used to access the resource.",
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
				Transform:   transform.From(extractLoadBalancerNameFromProbeID),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the probe resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProbePropertiesFormat.ProvisioningState"),
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
				Name:        "interval_in_seconds",
				Description: "The interval, in seconds, for how frequently to probe the endpoint for health status. Typically, the interval is slightly less than half the allocated timeout period (in seconds) which allows two full probes before taking the instance out of rotation. The default value is 15, the minimum value is 5.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ProbePropertiesFormat.IntervalInSeconds"),
			},
			{
				Name:        "number_of_probes",
				Description: "The number of probes where if no response, will result in stopping further traffic from being delivered to the endpoint. This values allows endpoints to be taken out of rotation faster or slower than the typical times used in Azure.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ProbePropertiesFormat.NumberOfProbes"),
			},
			{
				Name:        "port",
				Description: "The port for communicating the probe. Possible values range from 1 to 65535, inclusive.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ProbePropertiesFormat.Port"),
			},
			{
				Name:        "protocol",
				Description: "The protocol of the end point. If 'Tcp' is specified, a received ACK is required for the probe to be successful. If 'Http' or 'Https' is specified, a 200 OK response from the specifies URI is required for the probe to be successful. Possible values include: 'HTTP', 'TCP', 'HTTPS'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProbePropertiesFormat.Protocol"),
			},
			{
				Name:        "request_path",
				Description: "The URI used for requesting health status from the VM. Path is required if a protocol is set to http. Otherwise, it is not allowed. There is no default value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProbePropertiesFormat.RequestPath"),
			},
			{
				Name:        "load_balancing_rules",
				Description: "The load balancer rules that use this probe.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProbePropertiesFormat.LoadBalancingRules"),
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

func listLoadBalancerProbes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of load balancer
	loadBalancer := h.Item.(network.LoadBalancer)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*loadBalancer.ID, "/")[4]

	listLoadBalancerProbesClient := network.NewLoadBalancerProbesClient(subscriptionID)
	listLoadBalancerProbesClient.Authorizer = session.Authorizer

	result, err := listLoadBalancerProbesClient.List(ctx, resourceGroup, *loadBalancer.Name)
	if err != nil {
		return nil, err
	}
	for _, probe := range result.Values() {
		d.StreamListItem(ctx, probe)
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
		for _, probe := range result.Values() {
			d.StreamListItem(ctx, probe)
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

func getLoadBalancerProbe(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoadBalancerProbe")

	loadBalancerName := d.KeyColumnQuals["load_balancer_name"].GetStringValue()
	probeName := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty loadBalancerName, probeName or resourceGroup
	if loadBalancerName == "" || probeName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	LoadBalancerProbeClient := network.NewLoadBalancerProbesClient(subscriptionID)
	LoadBalancerProbeClient.Authorizer = session.Authorizer

	op, err := LoadBalancerProbeClient.Get(ctx, resourceGroup, loadBalancerName, probeName)
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

func extractLoadBalancerNameFromProbeID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(network.Probe)
	loadBalancerName := strings.Split(*data.ID, "/")[8]
	return loadBalancerName, nil
}
