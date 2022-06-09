package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/appplatform/mgmt/2020-07-01/appplatform"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureSpringCloudService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_spring_cloud_service",
		Description: "Azure Spring Cloud Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getSpringCloudService,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listSpringCloudServices,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource Id for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the Service. Possible values include: 'Creating', 'Updating', 'Deleting', 'Deleted', 'Succeeded', 'Failed', 'Moving', 'Moved', 'MoveFailed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_id",
				Description: "Service instance entity GUID which uniquely identifies a created resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ServiceID"),
			},
			{
				Name:        "sku_name",
				Description: "Name of the Sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "Tier of the Sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "sku_capacity",
				Description: "Current capacity of the target resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "version",
				Description: "Version of the service.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Version"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSpringCloudServiceDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "network_profile",
				Description: "Network profile of the service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSpringCloudServiceNetworkProfile),
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
				Transform:   transform.FromField("Tags"),
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

type SpringCloudServiceNetworkProfile struct {
	ServiceRuntimeSubnetID             *string
	AppSubnetID                        *string
	ServiceCidr                        *string
	ServiceRuntimeNetworkResourceGroup *string
	AppNetworkResourceGroup            *string
	OutboundPublicIPs                  *[]string
}

//// LIST FUNCTION

func listSpringCloudServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	// Get the details of the resource group
	resourceGroup := h.Item.(resources.Group)
	if resourceGroup.Name == nil {
		return nil, nil
	}

	client := appplatform.NewServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx, *resourceGroup.Name)
	if err != nil {
		plugin.Logger(ctx).Error("listSpringCloudServices", "list", err)
		return nil, err
	}

	for _, service := range result.Values() {
		d.StreamListItem(ctx, service)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listSpringCloudServices", "list_paging", err)
			return nil, err
		}
		for _, service := range result.Values() {
			d.StreamListItem(ctx, service)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSpringCloudService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSpringCloudService")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := appplatform.NewServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	service, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getSpringCloudService", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if service.ID != nil {
		return service, nil
	}

	return nil, nil
}

func listSpringCloudServiceDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSpringCloudServiceDiagnosticSettings")
	id := *h.Item.(appplatform.ServiceResource).ID

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
		plugin.Logger(ctx).Error("listSpringCloudServiceDiagnosticSettings", "list", err)
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

//// TRANSFORM FUNCTION

// If we return the API response directly, the output does not provide
// all the properties of NetworkProfile
func extractSpringCloudServiceNetworkProfile(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	workspace := d.HydrateItem.(appplatform.ServiceResource)
	var properties SpringCloudServiceNetworkProfile

	if workspace.Properties.NetworkProfile != nil {
		if workspace.Properties.NetworkProfile.ServiceRuntimeSubnetID != nil {
			properties.ServiceRuntimeSubnetID = workspace.Properties.NetworkProfile.ServiceRuntimeSubnetID
		}
		if workspace.Properties.NetworkProfile.AppSubnetID != nil {
			properties.AppSubnetID = workspace.Properties.NetworkProfile.AppSubnetID
		}
		if workspace.Properties.NetworkProfile.ServiceCidr != nil {
			properties.ServiceCidr = workspace.Properties.NetworkProfile.ServiceCidr
		}
		if workspace.Properties.NetworkProfile.ServiceRuntimeNetworkResourceGroup != nil {
			properties.ServiceRuntimeNetworkResourceGroup = workspace.Properties.NetworkProfile.ServiceRuntimeNetworkResourceGroup
		}
		if workspace.Properties.NetworkProfile.AppNetworkResourceGroup != nil {
			properties.AppNetworkResourceGroup = workspace.Properties.NetworkProfile.AppNetworkResourceGroup
		}
		if workspace.Properties.NetworkProfile.OutboundIPs != nil {
			if workspace.Properties.NetworkProfile.OutboundIPs.PublicIPs != nil {
				properties.OutboundPublicIPs = workspace.Properties.NetworkProfile.OutboundIPs.PublicIPs
			}
		}
	}

	return properties, nil
}
