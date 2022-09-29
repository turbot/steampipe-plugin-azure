package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2021-06-01-preview/eventgrid"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureEventGridTopic(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_eventgrid_topic",
		Description: "Azure Event Grid Topic",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getEventGridTopic,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "400", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listEventGridTopics,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified identifier of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the event grid topic resource. Possible values include: 'Creating', 'Updating', 'Deleting', 'Succeeded', 'Canceled', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TopicProperties.ProvisioningState"),
			},
			{
				Name:        "created_at",
				Description: "The timestamp of resource creation (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SystemData.CreatedAt").Transform(convertDateToTime),
			},
			{
				Name:        "created_by",
				Description: "The identity that created the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.CreatedBy"),
			},
			{
				Name:        "created_by_type",
				Description: "The type of identity that created the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.CreatedByType"),
			},
			{
				Name:        "disable_local_auth",
				Description: "This boolean is used to enable or disable local auth. Default value is false. When the property is set to true, only AAD token will be used to authenticate if user is allowed to publish to the topic.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("TopicProperties.DisableLocalAuth"),
				Default:     false,
			},
			{
				Name:        "endpoint",
				Description: "Endpoint for the event grid topic resource which is used for publishing the events.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TopicProperties.Endpoint"),
			},
			{
				Name:        "input_schema",
				Description: "This determines the format that event grid should expect for incoming events published to the event grid topic resource. Possible values include: 'EventGridSchema', 'CustomEventSchema', 'CloudEventSchemaV10'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TopicProperties.InputSchema"),
			},
			{
				Name:        "kind",
				Description: "Kind of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_at",
				Description: "The timestamp of resource last modification (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SystemData.LastModifiedAt").Transform(convertDateToTime),
			},
			{
				Name:        "last_modified_by",
				Description: "The identity that last modified the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.LastModifiedBy"),
			},
			{
				Name:        "last_modified_by_type",
				Description: "The type of identity that last modified the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.LastModifiedByType"),
			},
			{
				Name:        "location",
				Description: "Location of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_network_access",
				Description: "This determines if traffic is allowed over public network. By default it is enabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TopicProperties.PublicNetworkAccess"),
			},
			{
				Name:        "sku_name",
				Description: "Name of this SKU. Possible values include: 'Basic', 'Standard'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the eventgrid topic.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listEventGridTopicDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "extended_location",
				Description: "Extended location of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "Identity information for the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "inbound_ip_rules",
				Description: "This can be used to restrict traffic from specific IPs instead of all IPs. Note: These are considered only if PublicNetworkAccess is enabled.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TopicProperties.InboundIPRules"),
			},
			{
				Name:        "input_schema_mapping",
				Description: "Information about the InputSchemaMapping which specified the info about mapping event payload.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TopicProperties.InputSchemaMapping"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "List of private endpoint connections for the event grid topic.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractEventgridTopicPrivaterEndPointConnections),
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
				Transform:   transform.FromField("Location").Transform(formatRegion).Transform(toLower),
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

func listEventGridTopics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEventGridTopics")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := eventgrid.NewTopicsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListBySubscription(ctx, "", nil)
	if err != nil {
		plugin.Logger(ctx).Error("listEventGridTopics", "ListBySubscription", err)
		return nil, err
	}

	for _, topic := range result.Values() {
		d.StreamListItem(ctx, topic)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listEventGridTopics", "ListBySubscription_pagination", err)
			return nil, err
		}

		for _, topic := range result.Values() {
			d.StreamListItem(ctx, topic)
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEventGridTopic(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEventGridTopic")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := eventgrid.NewTopicsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getEventGridTopic", "get", err)
		return nil, err
	}

	return op, nil
}

func listEventGridTopicDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEventGridTopicDiagnosticSettings")
	id := *h.Item.(eventgrid.Topic).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Pagination is not supported
	op, err := client.List(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("listEventGridTopicDiagnosticSettings", "list", err)
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

//// TRANSFORM FUNCTIONS

// If we return the private endpoint connection directly from api response we will not receive all the properties of private endpoint connections.
func extractEventgridTopicPrivaterEndPointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("extractEventgridTopicPrivaterEndPointConnections")
	topic := d.HydrateItem.(eventgrid.Topic)
	var privateEndpointConnectionsInfo []map[string]interface{}
	if topic.PrivateEndpointConnections != nil {
		privateEndpointConnections := *topic.PrivateEndpointConnections
		for _, endpoint := range privateEndpointConnections {
			objectMap := make(map[string]interface{})

			if endpoint.ID != nil {
				objectMap["id"] = endpoint.ID
			}

			if endpoint.Name != nil {
				objectMap["name"] = endpoint.Name
			}

			if endpoint.Type != nil {
				objectMap["type"] = endpoint.Type
			}

			if endpoint.PrivateEndpointConnectionProperties != nil {
				if endpoint.PrivateEndpointConnectionProperties.PrivateEndpoint != nil {
					if endpoint.PrivateEndpointConnectionProperties.PrivateEndpoint.ID != nil {
						objectMap["endpointId"] = endpoint.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
					}
				}
				if endpoint.PrivateEndpointConnectionProperties.GroupIds != nil {
					objectMap["groupIds"] = endpoint.PrivateEndpointConnectionProperties.GroupIds
				}
				if endpoint.PrivateEndpointConnectionProperties.ProvisioningState != "" {
					objectMap["provisioningState"] = endpoint.PrivateEndpointConnectionProperties.ProvisioningState
				}
				if endpoint.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState != nil {
					if endpoint.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status != "" {
						objectMap["privateLinkServiceConnectionStateStatus"] = endpoint.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status
					}
					if endpoint.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description != nil {
						objectMap["privateLinkServiceConnectionStateDescription"] = endpoint.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description
					}
					if endpoint.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired != nil {
						objectMap["privateLinkServiceConnectionStateActionsRequired"] = endpoint.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired
					}
				}
			}
			privateEndpointConnectionsInfo = append(privateEndpointConnectionsInfo, objectMap)
		}
	}
	return privateEndpointConnectionsInfo, nil
}
