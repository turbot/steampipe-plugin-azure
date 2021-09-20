package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2019-03-01/servicefabric"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureServiceFabricCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_service_fabric_cluster",
		Description: "Azure Service Fabric Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getServiceFabricCluster,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceFabricClusters,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Azure resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Azure resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the cluster resource. Possible values include: 'Updating', 'Succeeded', 'Failed', 'Canceled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "Azure resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_code_version",
				Description: "The service fabric runtime version of the cluster. This property can only by set the user when **upgradeMode** is set to 'Manual'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ClusterCodeVersion"),
			},
			{
				Name:        "cluster_endpoint",
				Description: "The azure resource provider endpoint. A system service in the cluster connects to this  endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ClusterEndpoint"),
			},
			{
				Name:        "cluster_id",
				Description: "A service generated unique identifier for the cluster resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ClusterID"),
			},
			{
				Name:        "cluster_state",
				Description: "The current state of the cluster. Possible values include: 'WaitingForNodes', 'Deploying', 'BaselineUpgrade', 'UpdatingUserConfiguration', 'UpdatingUserCertificate', 'UpdatingInfrastructure', 'EnforcingClusterVersion', 'UpgradeServiceUnreachable', 'AutoScale', 'Ready'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ClusterState"),
			},
			{
				Name:        "event_store_service_enabled",
				Description: "Indicates if the event store service is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ClusterProperties.EventStoreServiceEnabled"),
				Default:     false,
			},
			{
				Name:        "etag",
				Description: "Azure resource etag.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "management_endpoint",
				Description: "The http management endpoint of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ManagementEndpoint"),
			},
			{
				Name:        "reliability_level",
				Description: "The reliability level sets the replica set size of system services. Possible values include: 'None', 'Bronze', 'Silver', 'Gold', 'Platinum'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ReliabilityLevel"),
			},
			{
				Name:        "upgrade_mode",
				Description: "The upgrade mode of the cluster when new service fabric runtime version is available. Possible values include: 'Automatic', 'Manual'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.UpgradeMode"),
			},
			{
				Name:        "vm_image",
				Description: "The VM image VMSS has been configured with. Generic names such as Windows or Linux can be used.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.VMImage"),
			},
			{
				Name:        "add_on_features",
				Description: "The list of add-on features to enable in the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.AddOnFeatures"),
			},
			{
				Name:        "available_cluster_versions",
				Description: "The service fabric runtime versions available for this cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.AvailableClusterVersions"),
			},
			{
				Name:        "azure_active_directory",
				Description: "The azure active directory authentication settings of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.AzureActiveDirectory"),
			},
			{
				Name:        "certificate",
				Description: "The certificate to use for securing the cluster. The certificate provided will be used for node to node security within the cluster, SSL certificate for cluster management endpoint and default admin client.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.Certificate"),
			},
			{
				Name:        "certificate_common_names",
				Description: "Describes a list of server certificates referenced by common name that are used to secure the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.CertificateCommonNames"),
			},
			{
				Name:        "client_certificate_common_names",
				Description: "The list of client certificates referenced by common name that are allowed to manage the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.ClientCertificateCommonNames"),
			},
			{
				Name:        "client_certificate_thumbprints",
				Description: "The list of client certificates referenced by thumbprint that are allowed to manage the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.ClientCertificateThumbprints"),
			},
			{
				Name:        "diagnostics_storage_account_config",
				Description: "The storage account information for storing service fabric diagnostic logs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.DiagnosticsStorageAccountConfig"),
			},
			{
				Name:        "fabric_settings",
				Description: "The list of custom fabric settings to configure the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.FabricSettings"),
			},
			{
				Name:        "node_types",
				Description: "The list of node types in the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.NodeTypes"),
			},
			{
				Name:        "reverse_proxy_certificate",
				Description: "The server certificate used by reverse proxy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.ReverseProxyCertificate"),
			},
			{
				Name:        "reverse_proxy_certificate_common_names",
				Description: "Describes a list of server certificates referenced by common name that are used to secure the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.ReverseProxyCertificateCommonNames"),
			},
			{
				Name:        "upgrade_description",
				Description: "The policy to use when upgrading the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.UpgradeDescription"),
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

func listServiceFabricClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	clusterClient := servicefabric.NewClustersClient(subscriptionID)
	clusterClient.Authorizer = session.Authorizer

	result, err := clusterClient.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listServiceFabricClusters", "list", err)
		return nil, err
	}

	// The API provides an URL for next set of data but accepts no param to implement pagination
	for _, cluster := range *result.Value {
		d.StreamListItem(ctx, cluster)
	}
	
	return nil, err
}

//// HYDRATE FUNCTIONS

func getServiceFabricCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServiceFabricCluster")

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

	clusterClient := servicefabric.NewClustersClient(subscriptionID)
	clusterClient.Authorizer = session.Authorizer

	cluster, err := clusterClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getServiceFabricCluster", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if cluster.ID != nil {
		return cluster, nil
	}

	return nil, nil
}
