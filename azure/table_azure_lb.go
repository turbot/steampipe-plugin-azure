package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lb",
		Description: "Azure Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getLoadBalancer,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLoadBalancers,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the load balancer resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "extended_location_name",
				Description: "The name of the extended location.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExtendedLocation.Name"),
			},
			{
				Name:        "extended_location_type",
				Description: "The type of the extended location. Possible values include: 'ExtendedLocationTypesEdgeZone'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExtendedLocation.Type").Transform(transform.ToString),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the load balancer resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "sku_name",
				Description: "Name of the load balancer SKU. Possible values include: 'Basic', 'Standard', 'Gateway'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Tier of the load balancer SKU. Possible values include: 'Regional', 'Global'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "backend_address_pools",
				Description: "Collection of backend address pools used by the load balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.BackendAddressPools"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listLoadBalancerDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "frontend_ip_configurations",
				Description: "Object representing the frontend IPs to be used for the load balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.FrontendIPConfigurations"),
			},
			{
				Name:        "inbound_nat_pools",
				Description: "Defines an external port range for inbound NAT to a single backend port on NICs associated with the load balancer. Inbound NAT rules are created automatically for each NIC associated with the Load Balancer using an external port from this range. Defining an Inbound NAT pool on the Load Balancer is mutually exclusive with defining inbound Nat rules. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an inbound NAT pool. They have to reference individual inbound NAT rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.InboundNatPools"),
			},
			{
				Name:        "inbound_nat_rules",
				Description: "Collection of inbound NAT Rules used by the load balancer. Defining inbound NAT rules on the load balancer is mutually exclusive with defining an inbound NAT pool. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an Inbound NAT pool. They have to reference individual inbound NAT rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.InboundNatRules"),
			},
			{
				Name:        "load_balancing_rules",
				Description: "Object collection representing the load balancing rules Gets the provisioning.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.LoadBalancingRules"),
			},
			{
				Name:        "outbound_rules",
				Description: "The outbound rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.OutboundRules"),
			},
			{
				Name:        "probes",
				Description: "Collection of probe objects used in the load balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.Probes"),
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
		}),
	}
}

//// LIST FUNCTIONS

func listLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	loadBalancersClient := network.NewLoadBalancersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	loadBalancersClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &loadBalancersClient, d.Connection)

	result, err := loadBalancersClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, loadBalancer := range result.Values() {
		d.StreamListItem(ctx, loadBalancer)
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
		for _, loadBalancer := range result.Values() {
			d.StreamListItem(ctx, loadBalancer)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getLoadBalancer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoadBalancer")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	loadBalancersClient := network.NewLoadBalancersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	loadBalancersClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &loadBalancersClient, d.Connection)

	op, err := loadBalancersClient.Get(ctx, resourceGroup, name, "")
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

func listLoadBalancerDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureLoadBalancerDiagnosticSettings")
	id := *h.Item.(network.LoadBalancer).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.List(ctx, id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["id"] = i.ID
		}
		if i.Name != nil {
			objectMap["name"] = i.Name
		}
		if i.Type != nil {
			objectMap["type"] = i.Type
		}
		if i.DiagnosticSettings != nil {
			objectMap["properties"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}
	return diagnosticSettings, nil
}
