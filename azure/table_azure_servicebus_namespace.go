package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2018-01-01-preview/servicebus"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureServiceBusNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_servicebus_namespace",
		Description: "Azure ServiceBus Namespace",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getServiceBusNamespace,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404", "400"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceBusNamespaces,
		},
		Columns: []*plugin.Column{
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

func listServiceBusNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listServiceBusNamespaces")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := servicebus.NewNamespacesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := client.List(context.Background())
		if err != nil {
			return nil, err
		}

		for _, namespace := range result.Values() {
			d.StreamListItem(ctx, namespace)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
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
	client := servicebus.NewNamespacesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(context.Background(), resourceGroup, name)
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
	client := servicebus.NewNamespacesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	data := h.Item.(servicebus.SBNamespace)
	resourceGroup := strings.Split(*data.ID, "/")[4]

	op, err := client.GetNetworkRuleSet(context.Background(), resourceGroup, *data.Name)
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

	client := insights.NewDiagnosticSettingsClient(subscriptionID)
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
