package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2021-06-01-preview/servicebus"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureServiceBusNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_servicebus_namespace",
		Description: "Azure ServiceBus Namespace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getServiceBusNamespace,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404", "400"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceBusNamespaces,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SBNamespaceProperties.ProvisioningState"),
			},
			{
				Name:        "zone_redundant",
				Description: "Enabling this property creates a Premium Service Bus Namespace in regions supported availability zones.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SBNamespaceProperties.ZoneRedundant"),
			},
			{
				Name:        "created_at",
				Description: "The time the namespace was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SBNamespaceProperties.CreatedAt").Transform(convertDateToTime),
			},
			{
				Name:        "metric_id",
				Description: "The identifier for Azure insights metrics.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SBNamespaceProperties.MetricID"),
			},
			{
				Name:        "servicebus_endpoint",
				Description: "Specifies the endpoint used to perform Service Bus operations.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SBNamespaceProperties.ServiceBusEndpoint"),
			},
			{
				Name:        "sku_capacity",
				Description: "The specified messaging units for the tier. For Premium tier, capacity are 1,2 and 4.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_name",
				Description: "Name of this SKU. Valid valuer are: 'Basic', 'Standard', 'Premium'.",
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
				Transform:   transform.FromField("SBNamespaceProperties.UpdatedAt").Transform(convertDateToTime),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the servicebus namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listServiceBusNamespaceDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "Specifies the properties of BYOK encryption configuration. Customer-managed key encryption at rest (Bring Your Own Key) is only available on Premium namespaces.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SBNamespaceProperties.Encryption"),
			},
			{
				Name:        "network_rule_set",
				Description: "Describes the network rule set for specified namespace. The ServiceBus Namespace must be Premium in order to attach a ServiceBus Namespace Network Rule Set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceBusNamespaceNetworkRuleSet,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The private endpoint connections of the namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listServiceBusNamespacePrivateEndpointConnections,
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

func listServiceBusNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listServiceBusNamespaces")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := servicebus.NewNamespacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, namespace := range result.Values() {
		d.StreamListItem(ctx, namespace)
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

		for _, namespace := range result.Values() {
			d.StreamListItem(ctx, namespace)
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

func getServiceBusNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceBusNamespace")

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
	client := servicebus.NewNamespacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getServiceBusNamespaceNetworkRuleSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceBusNamespaceNetworkRuleSet")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := servicebus.NewNamespacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	data := h.Item.(servicebus.SBNamespace)
	resourceGroup := strings.Split(*data.ID, "/")[4]

	op, err := client.GetNetworkRuleSet(ctx, resourceGroup, *data.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listServiceBusNamespaceDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listServiceBusNamespaceDiagnosticSettings")
	id := *h.Item.(servicebus.SBNamespace).ID

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

func listServiceBusNamespacePrivateEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listServiceBusNamespacePrivateEndpointConnections")

	namespace := h.Item.(servicebus.SBNamespace)
	resourceGroup := strings.Split(string(*namespace.ID), "/")[4]
	namespaceName := *namespace.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := servicebus.NewPrivateEndpointConnectionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, resourceGroup, namespaceName)
	if err != nil {
		plugin.Logger(ctx).Error("listServiceBusNamespacePrivateEndpointConnections", "list", err)
		return nil, err
	}

	var serviceBusNamespacePrivateEndpointConnections []map[string]interface{}

	for _, i := range op.Values() {
		serviceBusNamespacePrivateEndpointConnections = append(serviceBusNamespacePrivateEndpointConnections, extractServiceBusNamespacePrivateEndpointConnection(i))
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listServiceBusNamespacePrivateEndpointConnections", "list_paging", err)
			return nil, err
		}
		for _, i := range op.Values() {
			serviceBusNamespacePrivateEndpointConnections = append(serviceBusNamespacePrivateEndpointConnections, extractServiceBusNamespacePrivateEndpointConnection(i))
		}
	}

	return serviceBusNamespacePrivateEndpointConnections, nil
}

// If we return the API response directly, the output will not provide the properties of PrivateEndpointConnections
func extractServiceBusNamespacePrivateEndpointConnection(i servicebus.PrivateEndpointConnection) map[string]interface{} {
	serviceBusNamespacePrivateEndpointConnection := make(map[string]interface{})
	if i.ID != nil {
		serviceBusNamespacePrivateEndpointConnection["id"] = *i.ID
	}
	if i.Name != nil {
		serviceBusNamespacePrivateEndpointConnection["name"] = *i.Name
	}
	if i.Type != nil {
		serviceBusNamespacePrivateEndpointConnection["type"] = *i.Type
	}
	if i.PrivateEndpointConnectionProperties != nil {
		if len(i.PrivateEndpointConnectionProperties.ProvisioningState) > 0 {
			serviceBusNamespacePrivateEndpointConnection["provisioningState"] = i.PrivateEndpointConnectionProperties.ProvisioningState
		}
		if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState != nil {
			serviceBusNamespacePrivateEndpointConnection["privateLinkServiceConnectionState"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState
		}
		if i.PrivateEndpointConnectionProperties.PrivateEndpoint != nil && i.PrivateEndpointConnectionProperties.PrivateEndpoint.ID != nil {
			serviceBusNamespacePrivateEndpointConnection["privateEndpointPropertyID"] = i.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
		}
	}
	return serviceBusNamespacePrivateEndpointConnection
}
