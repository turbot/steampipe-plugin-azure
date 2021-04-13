package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-02-01/containerservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureKubernetesCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_kubernetes_cluster",
		Description: "Azure Kubernetes Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getKubernetesCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listKubernetesClusters,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the cluster.",
			},
			{
				Name:        "id",
				Description: "The ID of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location where the cluster is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku",
				Description: "The managed cluster SKU.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKubernetesCluster,
			},
			{
				Name:        "identity",
				Description: "The identity of the managed cluster, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "managed_cluster_properties",
				Description: "Properties of a managed cluster.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
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

//// FETCH FUNCTIONS ////

func listKubernetesClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := containerservice.NewManagedClustersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := client.List(ctx)
		if err != nil {
			return nil, err
		}

		for _, cluster := range result.Values() {
			d.StreamListItem(ctx, cluster)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getKubernetesCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKubernetesCluster")

	resourceName := ""
	resourceGroupName := ""
	if h.Item != nil {
		managedCluster := h.Item.(containerservice.ManagedCluster)
		resourceName = *managedCluster.Name
		resourceGroupName = strings.Split(string(*managedCluster.ID), "/")[4]
	} else {
		resourceName = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroupName = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := containerservice.NewManagedClustersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, resourceName)
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
