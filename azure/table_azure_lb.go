package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lb",
		Description: "Azure Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getLoadBalancer,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listLoadBalancers,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the load balancer resource. Possible values include: 'ProvisioningStateSucceeded', 'ProvisioningStateUpdating', 'ProvisioningStateDeleting', 'ProvisioningStateFailed'",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "Resource type.",
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
				Description: "Name of a load balancer SKU. Possible values include: 'LoadBalancerSkuNameBasic', 'LoadBalancerSkuNameStandard', 'LoadBalancerSkuNameGateway'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Tier of a load balancer SKU. Possible values include: 'LoadBalancerSkuTierRegional', 'LoadBalancerSkuTierGlobal'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "backend_address_pools",
				Description: "Collection of backend address pools used by a load balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.BackendAddressPools"),
			},
			{
				Name:        "frontend_ip_configurations",
				Description: "Object representing the frontend IPs to be used for the load balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.FrontendIPConfigurations"),
			},
			{
				Name:        "inbound_nat_pools",
				Description: "Defines an external port range for inbound NAT to a single backend port on NICs associated with a load balancer. Inbound NAT rules are created automatically for each NIC associated with the Load Balancer using an external port from this range. Defining an Inbound NAT pool on your Load Balancer is mutually exclusive with defining inbound Nat rules. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an inbound NAT pool. They have to reference individual inbound NAT rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LoadBalancerPropertiesFormat.InboundNatPools"),
			},
			{
				Name:        "inbound_nat_rules",
				Description: "Collection of inbound NAT Rules used by a load balancer. Defining inbound NAT rules on your load balancer is mutually exclusive with defining an inbound NAT pool. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an Inbound NAT pool. They have to reference individual inbound NAT rules.",
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

			// Azure standard column
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

//// FETCH FUNCTIONS ////

func listLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	LoadBalancersClient := network.NewLoadBalancersClient(subscriptionID)
	LoadBalancersClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := LoadBalancersClient.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		for _, loadBalancer := range result.Values() {
			d.StreamListItem(ctx, loadBalancer)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getLoadBalancer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoadBalancer")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	LoadBalancersClient := network.NewLoadBalancersClient(subscriptionID)
	LoadBalancersClient.Authorizer = session.Authorizer

	op, err := LoadBalancersClient.Get(ctx, resourceGroup, name, "")
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
