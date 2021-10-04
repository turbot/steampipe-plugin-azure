package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/hybridkubernetes/mgmt/hybridkubernetes"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/kubernetesconfiguration/mgmt/kubernetesconfiguration"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureHybridKubernetesConnectedCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_hybrid_kubernetes_connected_cluster",
		Description: "Azure Hybrid Kubernetes Connected Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getArcKubernetesCluster,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listArcKubernetesClusters,
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
				Name:        "connectivity_status",
				Description: "Represents the connectivity status of the connected cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.ConnectivityStatus"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the connected cluster resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_public_key_certificate",
				Description: "Base64 encoded public certificate used by the agent to do the initial handshake to the backend services in Azure.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.AgentPublicKeyCertificate"),
			},
			{
				Name:        "agent_version",
				Description: "Version of the agent running on the connected cluster resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.AgentVersion"),
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
				Name:        "distribution",
				Description: "The Kubernetes distribution running on this connected cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.Distribution"),
			},
			{
				Name:        "infrastructure",
				Description: "The infrastructure on which the Kubernetes cluster represented by this connected cluster is running on.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.Infrastructure"),
			},
			{
				Name:        "kubernetes_version",
				Description: "The Kubernetes version of the connected cluster resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.KubernetesVersion"),
			},
			{
				Name:        "last_connectivity_time",
				Description: "Time representing the last instance when heart beat was received from the cluster.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ConnectedClusterProperties.LastConnectivityTime").Transform(convertDateToTime),
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
				Name:        "managed_identity_certificate_expiration_time",
				Description: "Expiration time of the managed identity certificate.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ConnectedClusterProperties.ManagedIdentityCertificateExpirationTime").Transform(convertDateToTime),
			},
			{
				Name:        "offering",
				Description: "Connected cluster offering.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectedClusterProperties.Offering"),
			},
			{
				Name:        "total_core_count",
				Description: "Number of CPU cores present in the connected cluster resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ConnectedClusterProperties.TotalCoreCount"),
			},
			{
				Name:        "total_node_count",
				Description: "Number of nodes present in the connected cluster resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ConnectedClusterProperties.TotalNodeCount"),
			},
			{
				Name:        "extensions",
				Description: "The extensions of the connected cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listArcKubernetesClusterExtensions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity",
				Description: "The identity of the connected cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConnectedClusterProperties.Identity"),
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

func listArcKubernetesClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := hybridkubernetes.NewConnectedClusterClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listArcKubernetesClusters", "list", err)
		return nil, err
	}

	for _, cluster := range result.Values() {
		d.StreamListItem(ctx, cluster)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listArcKubernetesClusters", "list_paging", err)
			return nil, err
		}
		for _, cluster := range result.Values() {
			d.StreamListItem(ctx, cluster)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getArcKubernetesCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getArcKubernetesCluster")

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

	client := hybridkubernetes.NewConnectedClusterClient(subscriptionID)
	client.Authorizer = session.Authorizer

	cluster, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getArcKubernetesCluster", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if cluster.ID != nil {
		return cluster, nil
	}

	return nil, nil
}

func listArcKubernetesClusterExtensions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(hybridkubernetes.ConnectedCluster)
	resourceGroup := strings.Split(*cluster.ID, "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := kubernetesconfiguration.NewExtensionsClient(subscriptionID)
	client.Authorizer = session.Authorizer
	extensions := []kubernetesconfiguration.ExtensionInstance{}
	result, err := client.List(ctx, resourceGroup, "Microsoft.Kubernetes", "connectedClusters", *cluster.Name)
	if err != nil {
		plugin.Logger(ctx).Error("listArcKubernetesClusterExtensions", "list", err)
		return nil, err
	}

	extensions = append(extensions, result.Values()...)

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listArcKubernetesClusterExtensions", "list_paging", err)
			return nil, err
		}
		extensions = append(extensions, result.Values()...)
	}

	return extensions, nil
}
