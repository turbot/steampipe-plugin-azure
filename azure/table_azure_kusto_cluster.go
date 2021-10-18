package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2021-01-01/kusto"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureKustoCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_kusto_cluster",
		Description: "Azure Kusto Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getKustoCluster,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listKustoClusters,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource.",
			},
			{
				Name:        "id",
				Description: "The resource Id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioned state of the resource. Possible values include: 'Running', 'Creating', 'Deleting', 'Succeeded', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ProvisioningState"),
			},
			{
				Name:        "state",
				Description: "The state of the resource. Possible values include: 'Creating', 'Deleted', 'Deleting', 'Running', 'Starting', 'Stopped', 'Stopping', 'Unavailable'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.State"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "Specifies the name of the region, the resource is created at.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_ingestion_uri",
				Description: "The cluster data ingestion URI.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.DataIngestionURI"),
			},
			{
				Name:        "etag",
				Description: "An ETag of the resource created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_disk_encryption",
				Description: "A boolean value that indicates if the cluster's disks are encrypted.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ClusterProperties.EnableDiskEncryption"),
			},
			{
				Name:        "enable_double_encryption",
				Description: "A boolean value that indicates if double encryption is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ClusterProperties.EnableDoubleEncryption"),
			},
			{
				Name:        "enable_purge",
				Description: "A boolean value that indicates if the purge operations are enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ClusterProperties.EnablePurge"),
			},
			{
				Name:        "enable_streaming_ingest",
				Description: "A boolean value that indicates if the streaming ingest is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ClusterProperties.EnableStreamingIngest"),
			},
			{
				Name:        "engine_type",
				Description: "The engine type. Possible values include: 'EngineTypeV2', 'EngineTypeV3'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.EngineType"),
			},
			{
				Name:        "sku_capacity",
				Description: "SKU capacity of the resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_name",
				Description: "SKU name of the resource. Possible values include: 'KC8', 'KC16', 'KS8', 'KS16', 'D13V2', 'D14V2', 'L8', 'L16'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "SKU tier of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "state_reason",
				Description: "SKU tier of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SClusterPropertiesku.StateReason"),
			},
			{
				Name:        "uri",
				Description: "The cluster URI.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.URI"),
			},
			{
				Name:        "language_extensions",
				Description: "List of the cluster's language extensions.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.LanguageExtensions"),
			},
			{
				Name:        "key_vault_properties",
				Description: "KeyVault properties for the cluster encryption.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.KeyVaultProperties"),
			},
			{
				Name:        "optimized_autoscale",
				Description: "Optimized auto scale definition.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.OptimizedAutoscale"),
			},
			{
				Name:        "trusted_external_tenants",
				Description: "The cluster's external tenants.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.TrustedExternalTenants"),
			},
			{
				Name:        "virtual_network_configuration",
				Description: "Virtual network definition of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.VirtualNetworkConfiguration"),
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

func listKustoClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	kustoClient := kusto.NewClustersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	kustoClient.Authorizer = session.Authorizer

	//Pagination does not support for kusto cluster list call till date
	result, err := kustoClient.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listKustoClusters", "list", err)
		return nil, err
	}

	for _, cluster := range *result.Value {
		d.StreamListItem(ctx, cluster)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKustoCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKustoCluster")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	kustoClient := kusto.NewClustersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	kustoClient.Authorizer = session.Authorizer

	op, err := kustoClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getKustoCluster", "get", err)
		return nil, err
	}

	return op, nil
}
