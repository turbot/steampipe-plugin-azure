package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureEventHubNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_eventhub_namespace",
		Description: "Azure Event Hub Namespace",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAzureEventHubNamespace,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureEventHubNamespaces,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource Id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The Resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the Namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.ProvisioningState"),
			},
			{
				Name:        "created_at",
				Description: "The time the Namespace was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EHNamespaceProperties.CreatedAt"),
			},
			{
				Name:        "cluster_arm_id",
				Description: "Cluster ARM ID of the Namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.ClusterArmId"),
			},
			{
				Name:        "is_auto_inflate_enabled",
				Description: "Value that indicates whether AutoInflate is enabled for eventhub namespace.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EHNamespaceProperties.IsAutoInflateEnabled"),
			},
			{
				Name:        "kafka_enabled",
				Description: "Value that indicates whether Kafka is enabled for eventhub namespace.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EHNamespaceProperties.KafkaEnabled"),
			},
			{
				Name:        "maximum_throughput_units",
				Description: "Upper limit of throughput units when AutoInflate is enabled, value should be within 0 to 20 throughput units.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("EHNamespaceProperties.MaximumThroughputUnits"),
			},
			{
				Name:        "metric_id",
				Description: "Identifier for Azure Insights metrics.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.Metric_id"),
			},
			{
				Name:        "principal_id",
				Description: "ObjectId from the KeyVault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.PrincipalId"),
			},
			{
				Name:        "service_bus_endpoint",
				Description: "Endpoint you can use to perform Service Bus operations.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EHNamespaceProperties.ServiceBusEndpoint"),
			},
			{
				Name:        "tenant_id",
				Description: "TenantId from the KeyVault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.TenantId"),
			},
			{
				Name:        "updated_at",
				Description: "The time the Namespace was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EHNamespaceProperties.UpdatedAt"),
			},
			{
				Name:        "zone_redundant",
				Description: "Enabling this property creates a Standard Event Hubs Namespace in regions supported availability zones.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EHNamespaceProperties.ZoneRedundant"),
			},
			{
				Name:        "encryption",
				Description: "Cluster ARM ID of the Namespace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EHNamespaceProperties.Encryption"),
			},
			{
				Name:        "network_rule_set",
				Description: "Gets NetworkRuleSet for a Namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkRuleSet,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "sku",
				Description: "Properties of sku resource.",
				Type:        proto.ColumnType_JSON,
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

//// LIST FUNCTION ////

func listAzureEventHubNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureEventHubNamespaces")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := eventhub.NewNamespacesClient(subscriptionID)
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

//// HYDRATE FUNCTIONS ////

func getAzureEventHubNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeSnapshot")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := eventhub.NewNamespacesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(context.Background(), resourceGroup, name)
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

func getNetworkRuleSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getNetworkRuleSet")

	namespace:= h.Item.(eventhub.EHNamespace)
	resourceGroupName := strings.Split(string(*namespace.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	networkClient := eventhub.NewNamespacesClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	op, err := networkClient.GetNetworkRuleSet(context.Background(), resourceGroupName, *namespace.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
