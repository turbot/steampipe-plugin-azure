package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/azurestackhci/mgmt/azurestackhci"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureArcCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_arc_Cluster",
		Description: "Azure Arc Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getArcCluster,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listArcClusters,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the configuration store. Possible values include: 'Creating', 'Updating', 'Deleting', 'Succeeded', 'Failed', 'Canceled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigurationStoreProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The creation date of configuration store.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ConfigurationStoreProperties.CreationDate").Transform(convertDateToTime),
			},
			{
				Name:        "endpoint",
				Description: "The DNS endpoint where the configuration store API will be available.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigurationStoreProperties.Endpoint"),
			},
			{
				Name:        "public_network_access",
				Description: "Control permission for data plane traffic coming from public networks while private endpoint is enabled. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getPublicNetworkAccess,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "sku_name",
				Description: "The SKU name of the configuration store.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the configuration store.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAppConfigurationDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "The encryption settings of the configuration store.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConfigurationStoreProperties.Encryption"),
			},
			{
				Name:        "identity",
				Description: "The managed identity information, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The list of private endpoint connections that are set up for this resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractAppConfigurationPrivateEndpointConnections),
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

func listArcClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := azurestackhci.NewClustersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listArcClusters", "list", err)
		return nil, err
	}

	for _, cluster := range result.Values() {
		d.StreamListItem(ctx, cluster)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listArcClusters", "list_paging", err)
			return nil, err
		}
		for _, cluster := range result.Values() {
			d.StreamListItem(ctx, cluster)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getArcCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getArcCluster")

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

	client := azurestackhci.NewClustersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	cluster, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getArcCluster", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if cluster.ID != nil {
		return cluster, nil
	}

	return nil, nil
}
