package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/frontdoor/mgmt/frontdoor"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureFrontDoor(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_frontdoor",
		Description: "Azure Front Door",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getFrontDoor,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listFrontDoors,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the front door.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cname",
				Description: "The host that each frontendEndpoint must CNAME to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Cname"),
			},
			{
				Name:        "enabled_state",
				Description: "Operational status of the front door load balancer. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EnabledState"),
			},
			{
				Name:        "friendly_name",
				Description: "A friendly name for the front door.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.FriendlyName"),
			},
			{
				Name:        "front_door_id",
				Description: "The ID of the front door.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.FrontdoorID"),
			},
			{
				Name:        "resource_state",
				Description: "Resource status of the front door. Possible values include: 'Creating', 'Enabling', 'Enabled', 'Disabling', 'Disabled', 'Deleting'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ResourceState"),
			},
			{
				Name:        "backend_pools",
				Description: "Backend pools available to routing rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.BackendPools"),
			},
			{
				Name:        "backend_pools_settings",
				Description: "Settings for all backend pools",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.BackendPoolsSettings"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listFrontDoorDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "frontend_endpoints",
				Description: "Frontend endpoints available to routing rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.FrontendEndpoints"),
			},
			{
				Name:        "health_probe_settings",
				Description: "Health probe settings associated with this Front Door instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.HealthProbeSettings"),
			},
			{
				Name:        "load_balancing_settings",
				Description: "Load balancing settings associated with this front door instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.LoadBalancingSettings"),
			},
			{
				Name:        "rules_engines",
				Description: "Rules engine configurations available to routing rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.RulesEngines"),
			},
			{
				Name:        "routing_rules",
				Description: "Routing rules associated with this front door.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.RoutingRules"),
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

//// LIST FUNCTION

func listFrontDoors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := frontdoor.NewFrontDoorsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listFrontDoors", "list", err)
		return nil, err
	}

	for _, door := range result.Values() {
		d.StreamListItem(ctx, door)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listFrontDoors", "list_paging", err)
			return nil, err
		}
		for _, door := range result.Values() {
			d.StreamListItem(ctx, door)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getFrontDoor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFrontDoor")

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

	client := frontdoor.NewFrontDoorsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	door, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getFrontDoor", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if door.ID != nil {
		return door, nil
	}

	return nil, nil
}

func listFrontDoorDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listFrontDoorDiagnosticSettings")
	id := *h.Item.(frontdoor.FrontDoor).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("listFrontDoorDiagnosticSettings", "list", err)
		return nil, err
	}

	// If we return the API response directly, the output does not provide
	// all the contents of DiagnosticSettings
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
