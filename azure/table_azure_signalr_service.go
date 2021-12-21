package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/signalr/mgmt/2020-05-01/signalr"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSignalRService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_signalr_service",
		Description: "Azure SignalR Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getSignalRService,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listSignalRServices,
		},
		Columns: []*plugin.Column{
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
				Description: "Provisioning state of the resource. Possible values include: 'Unknown', 'Succeeded', 'Failed', 'Canceled', 'Running', 'Creating', 'Updating', 'Deleting', 'Moving'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "external_ip",
				Description: "The publicly accessible IP of the SignalR service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ExternalIP"),
			},
			{
				Name:        "host_name",
				Description: "FQDN of the SignalR service instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HostName"),
			},
			{
				Name:        "host_name_prefix",
				Description: "Prefix for the host name of the SignalR service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HostNamePrefix"),
			},
			{
				Name:        "kind",
				Description: "The kind of the service. Possible values include: 'SignalR', 'RawWebSockets'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_port",
				Description: "The publicly accessible port of the SignalR service which is designed for browser/client side usage.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.PublicPort"),
			},
			{
				Name:        "server_port",
				Description: "The publicly accessible port of the SignalR service which is designed for customer server side usage.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.ServerPort"),
			},
			{
				Name:        "version",
				Description: "Version of the SignalR resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Version"),
			},
			{
				Name:        "cors",
				Description: "Cross-Origin Resource Sharing (CORS) settings of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Cors"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the SignalR service.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSignalRServiceDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "features",
				Description: "List of SignalR feature flags.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Features"),
			},
			{
				Name:        "network_acls",
				Description: "Network ACLs of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.NetworkACLs"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "Private endpoint connections to the SignalR resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSignalRServicePrivateEndpointConnections),
			},
			{
				Name:        "sku",
				Description: "The billing information of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "upstream",
				Description: "Upstream settings when the Azure SignalR is in server-less mode.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Upstream"),
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

type SignalRServicePrivateEndpointConnections struct {
	PrivateEndpointPropertyID         interface{}
	PrivateLinkServiceConnectionState interface{}
	ProvisioningState                 interface{}
	ID                                *string
	Name                              *string
	Type                              *string
}

//// LIST FUNCTION

func listSignalRServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := signalr.NewClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listSignalRServices", "list", err)
		return nil, err
	}

	for _, service := range result.Values() {
		d.StreamListItem(ctx, service)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listSignalRServices", "list_paging", err)
			return nil, err
		}
		for _, service := range result.Values() {
			d.StreamListItem(ctx, service)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSignalRService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSignalRService")

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

	client := signalr.NewClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getSignalRService", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func listSignalRServiceDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSignalRServiceDiagnosticSettings")
	id := *h.Item.(signalr.ResourceType).ID

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
		plugin.Logger(ctx).Error("listSignalRServiceDiagnosticSettings", "list", err)
		return nil, err
	}

	// If we return the API response directly, the output does not provide all
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

//// TRANSFORM FUNCTION

// If we return the API response directly, the output will not provide all the properties of PrivateEndpointConnections
func extractSignalRServicePrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	service := d.HydrateItem.(signalr.ResourceType)
	info := []SignalRServicePrivateEndpointConnections{}

	if service.Properties != nil && service.Properties.PrivateEndpointConnections != nil {
		for _, connection := range *service.Properties.PrivateEndpointConnections {
			properties := SignalRServicePrivateEndpointConnections{}
			properties.ID = connection.ID
			properties.Name = connection.Name
			properties.Type = connection.Type
			if connection.PrivateEndpointConnectionProperties != nil {
				if connection.PrivateEndpointConnectionProperties.PrivateEndpoint != nil {
					properties.PrivateEndpointPropertyID = connection.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
				}
				properties.PrivateLinkServiceConnectionState = connection.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState
				properties.ProvisioningState = connection.PrivateEndpointConnectionProperties.ProvisioningState
			}
			info = append(info, properties)
		}
	}

	return info, nil
}
