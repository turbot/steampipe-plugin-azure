package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/hdinsight/mgmt/hdinsight"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureHDInsightCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_hdinsight_cluster",
		Description: "Azure HDInsight Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getHDInsightCluster,
			Tags: map[string]string{
				"service": "Microsoft.HDInsight",
				"action":  "clusters/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listHDInsightClusters,
			Tags: map[string]string{
				"service": "Microsoft.HDInsight",
				"action":  "clusters/read",
			},
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
				Description: "The provisioning state, which only appears in the response. Possible values include: 'InProgress', 'Failed', 'Succeeded', 'Canceled', 'Deleting'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_hdp_version",
				Description: "The hdp version of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ClusterHdpVersion"),
			},
			{
				Name:        "cluster_id",
				Description: "The cluster id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ClusterID"),
			},
			{
				Name:        "cluster_state",
				Description: "The state of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ClusterState"),
			},
			{
				Name:        "cluster_version",
				Description: "The version of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ClusterVersion"),
			},
			{
				Name:        "created_date",
				Description: "The date on which the cluster was created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CreatedDate"),
			},
			{
				Name:        "etag",
				Description: "The ETag for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "min_supported_tls_version",
				Description: "The minimal supported tls version of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinSupportedTLSVersion"),
			},
			{
				Name:        "os_type",
				Description: "The type of operating system. Possible values include: 'Windows', 'Linux'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.OsType").Transform(transform.ToString),
			},
			{
				Name:        "tier",
				Description: "The cluster tier. Possible values include: 'Standard', 'Premium'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Tier").Transform(transform.ToString),
			},
			{
				Name:        "cluster_definition",
				Description: "The cluster definition.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ClusterDefinition"),
			},
			{
				Name:        "compute_isolation_properties",
				Description: "The compute isolation properties of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ComputeIsolationProperties"),
			},
			{
				Name:        "compute_profile",
				Description: "The complete profile of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ComputeProfile"),
			},
			{
				Name:        "connectivity_endpoints",
				Description: "The list of connectivity endpoints.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ConnectivityEndpoints"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listHDInsightClusterDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "disk_encryption_properties",
				Description: "The disk encryption properties of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.DiskEncryptionProperties"),
			},
			{
				Name:        "encryption_in_transit_properties",
				Description: "The encryption-in-transit properties of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.EncryptionInTransitProperties"),
			},
			{
				Name:        "errors",
				Description: "The list of errors.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Errors"),
			},
			{
				Name:        "excluded_services_config",
				Description: "The excluded services config of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ExcludedServicesConfig"),
			},
			{
				Name:        "identity",
				Description: "The identity of the cluster, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "kafka_rest_properties",
				Description: "The cluster kafka rest proxy configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.KafkaRestProperties"),
			},
			{
				Name:        "network_properties",
				Description: "The network properties of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.NetworkProperties"),
			},
			{
				Name:        "quota_info",
				Description: "The quota information of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.QuotaInfo"),
			},
			{
				Name:        "security_profile",
				Description: "The security profile of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.SecurityProfile"),
			},
			{
				Name:        "storage_profile",
				Description: "The storage profile of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.StorageProfile"),
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

func listHDInsightClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := hdinsight.NewClustersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listHDInsightClusters", "list", err)
		return nil, err
	}

	for _, cluster := range result.Values() {
		d.StreamListItem(ctx, cluster)
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listHDInsightClusters", "list_paging", err)
			return nil, err
		}
		for _, cluster := range result.Values() {
			d.StreamListItem(ctx, cluster)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getHDInsightCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHDInsightCluster")

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

	client := hdinsight.NewClustersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getHDInsightCluster", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func listHDInsightClusterDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listHDInsightClusterDiagnosticSettings")
	id := *h.Item.(hdinsight.Cluster).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.List(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("listHDInsightClusterDiagnosticSettings", "list", err)
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
