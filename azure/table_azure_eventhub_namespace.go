package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureEventHubNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_eventhub_namespace",
		Description: "Azure Event Hub Namespace",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getEventHubNamespace,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "400", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listEventHubNamespaces,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the resource.",
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
				Description: "Provisioning state of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.ProvisioningState"),
			},
			{
				Name:        "created_at",
				Description: "The time the namespace was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EHNamespaceProperties.CreatedAt").Transform(convertDateToTime),
			},
			{
				Name:        "cluster_arm_id",
				Description: "Cluster ARM ID of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.ClusterArmID"),
			},
			{
				Name:        "is_auto_inflate_enabled",
				Description: "Indicates whether auto-inflate is enabled for eventhub namespace.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EHNamespaceProperties.IsAutoInflateEnabled"),
			},
			{
				Name:        "kafka_enabled",
				Description: "Indicates whether kafka is enabled for eventhub namespace, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EHNamespaceProperties.KafkaEnabled"),
			},
			{
				Name:        "maximum_throughput_units",
				Description: "Upper limit of throughput units when auto-inflate is enabled, value should be within 0 to 20 throughput units.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("EHNamespaceProperties.MaximumThroughputUnits"),
			},
			{
				Name:        "metric_id",
				Description: "Identifier for azure insights metrics.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.Metric_id"),
			},
			{
				Name:        "service_bus_endpoint",
				Description: "Endpoint you can use to perform service bus operations.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.ServiceBusEndpoint"),
			},
			{
				Name:        "sku_capacity",
				Description: "The Event Hubs throughput units, value should be 0 to 20 throughput units.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_name",
				Description: "Name of this SKU. Possible values include: 'Basic', 'Standard'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "The billing tier of this particular SKU. Valid values are: 'Basic', 'Standard', 'Premium'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "updated_at",
				Description: "The time the namespace was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EHNamespaceProperties.UpdatedAt").Transform(convertDateToTime),
			},
			{
				Name:        "zone_redundant",
				Description: "Enabling this property creates a standard event hubs namespace in regions supported availability zones.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EHNamespaceProperties.ZoneRedundant"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the eventhub namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listEventHubNamespaceDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "Properties of BYOK encryption description.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EHNamespaceProperties.Encryption"),
			},
			{
				Name:        "identity",
				Description: "Describes the properties of BYOK encryption description.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EHNamespaceProperties.Encryption"),
			},
			{
				Name:        "network_rule_set",
				Description: "Describes the network rule set for specified namespace. The EventHub Namespace must be Premium in order to attach a EventHub Namespace Network Rule Set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkRuleSet,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The private endpoint connections of the namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listEventHubNamespacePrivateEndpointConnections,
				Transform:   transform.FromValue(),
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
				Transform:   transform.FromField("Location").Transform(formatRegion).Transform(toLower),
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

func listEventHubNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEventHubNamespaces")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := eventhub.NewNamespacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, namespace := range result.Values() {
		d.StreamListItem(ctx, namespace)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, namespace := range result.Values() {
			d.StreamListItem(ctx, namespace)
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEventHubNamespace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEventHubNamespace")

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
	client := eventhub.NewNamespacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getNetworkRuleSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getNetworkRuleSet")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	networkClient := eventhub.NewNamespacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	namespace := h.Item.(eventhub.EHNamespace)
	resourceGroupName := strings.Split(string(*namespace.ID), "/")[4]

	op, err := networkClient.GetNetworkRuleSet(ctx, resourceGroupName, *namespace.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listEventHubNamespaceDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEventHubNamespaceDiagnosticSettings")
	id := *h.Item.(eventhub.EHNamespace).ID

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

func listEventHubNamespacePrivateEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listEventHubNamespacePrivateEndpointConnections")

	namespace := h.Item.(eventhub.EHNamespace)
	resourceGroup := strings.Split(string(*namespace.ID), "/")[4]
	namespaceName := *namespace.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := eventhub.NewPrivateEndpointConnectionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, resourceGroup, namespaceName)
	if err != nil {
		plugin.Logger(ctx).Error("listEventHubNamespacePrivateEndpointConnections", "list", err)
		return nil, err
	}

	var eventHubNamespacePrivateEndpointConnections []map[string]interface{}

	for _, i := range op.Values() {
		eventHubNamespacePrivateEndpointConnections = append(eventHubNamespacePrivateEndpointConnections, extractEventHubNamespacePrivateEndpointConnections(i))
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listEventHubNamespacePrivateEndpointConnections", "list_paging", err)
			return nil, err
		}
		for _, i := range op.Values() {
			eventHubNamespacePrivateEndpointConnections = append(eventHubNamespacePrivateEndpointConnections, extractEventHubNamespacePrivateEndpointConnections(i))
		}
	}

	return eventHubNamespacePrivateEndpointConnections, nil
}

// If we return the API response directly, the output will not provide the properties of PrivateEndpointConnections

func extractEventHubNamespacePrivateEndpointConnections(i eventhub.PrivateEndpointConnection) map[string]interface{} {
	eventHubNamespacePrivateEndpointConnection := make(map[string]interface{})
	if i.ID != nil {
		eventHubNamespacePrivateEndpointConnection["id"] = *i.ID
	}
	if i.Name != nil {
		eventHubNamespacePrivateEndpointConnection["name"] = *i.Name
	}
	if i.Type != nil {
		eventHubNamespacePrivateEndpointConnection["type"] = *i.Type
	}
	if i.PrivateEndpointConnectionProperties != nil {
		if len(i.PrivateEndpointConnectionProperties.ProvisioningState) > 0 {
			eventHubNamespacePrivateEndpointConnection["provisioningState"] = i.PrivateEndpointConnectionProperties.ProvisioningState
		}
		if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState != nil {
			eventHubNamespacePrivateEndpointConnection["privateLinkServiceConnectionState"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState
		}
		if i.PrivateEndpointConnectionProperties.PrivateEndpoint != nil && i.PrivateEndpointConnectionProperties.PrivateEndpoint.ID != nil {
			eventHubNamespacePrivateEndpointConnection["privateEndpointPropertyID"] = i.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
		}
	}
	return eventHubNamespacePrivateEndpointConnection
}
