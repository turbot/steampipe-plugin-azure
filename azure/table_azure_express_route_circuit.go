package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureExpressRouteCircuit(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_express_route_circuit",
		Description: "Azure Express Route Circuit",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getExpressRouteCircuit,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listExpressRouteCircuits,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the circuit.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku_name",
				Description: "The name of the SKU.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "The tier of the SKU. Possible values include: 'Standard', 'Premium', 'Basic', 'Local'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "sku_family",
				Description: "The family of the SKU. Possible values include: 'UnlimitedData', 'MeteredData'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Family").Transform(transform.ToString),
			},
			{
				Name:        "allow_classic_operations",
				Description: "Allow classic operations.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.AllowClassicOperations"),
			},
			{
				Name:        "circuit_provisioning_state",
				Description: "The CircuitProvisioningState state of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.CircuitProvisioningState"),
			},
			{
				Name:        "service_provider_provisioning_state",
				Description: "The ServiceProviderProvisioningState state of the resource. Possible values include: 'NotProvisioned', 'Provisioning', 'Provisioned', 'Deprovisioning'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.ServiceProviderProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "authorizations",
				Description: "The list of authorizations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.Authorizations"),
			},
			{
				Name:        "peerings",
				Description: "The list of peerings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.Peerings"),
			},
			{
				Name:        "service_key",
				Description: "The ServiceKey.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.ServiceKey"),
			},
			{
				Name:        "service_provider_notes",
				Description: "The ServiceProviderNotes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.ServiceProviderNotes"),
			},
			{
				Name:        "service_provider_properties",
				Description: "The ServiceProviderProperties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties"),
			},
			{
				Name:        "express_route_port",
				Description: "The reference to the ExpressRoutePort resource when the circuit is provisioned on an ExpressRoutePort resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.ExpressRoutePort"),
			},
			{
				Name:        "bandwidth_in_gbps",
				Description: "The bandwidth of the circuit when the circuit is provisioned on an ExpressRoutePort resource.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.BandwidthInGbps"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the express route circuit resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "global_reach_enabled",
				Description: "Flag denoting global reach status.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ExpressRouteCircuitPropertiesFormat.GlobalReachEnabled"),
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
				Name:        "environment_name",
				Description: ColumnDescriptionEnvironmentName,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getEnvironmentName).WithCache(),
				Transform:   transform.FromValue(),
			},
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

//// LIST FUNCTION

func listExpressRouteCircuits(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	expressRouteCircuitClient := network.NewExpressRouteCircuitsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	expressRouteCircuitClient.Authorizer = session.Authorizer

	result, err := expressRouteCircuitClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, routeCircuit := range result.Values() {
		d.StreamListItem(ctx, routeCircuit)
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
		for _, routeCircuit := range result.Values() {
			d.StreamListItem(ctx, routeCircuit)
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

func getExpressRouteCircuit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getExpressRouteCircuit")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	expressRouteCircuitClient := network.NewExpressRouteCircuitsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	expressRouteCircuitClient.Authorizer = session.Authorizer

	op, err := expressRouteCircuitClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}
	return op, nil
}
