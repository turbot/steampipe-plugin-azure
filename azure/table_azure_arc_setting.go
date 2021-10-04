package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/azurestackhci/mgmt/azurestackhci"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureArcSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_arc_setting",
		Description: "Azure Arc Setting",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getArcSetting,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listArcClusters,
			Hydrate:       listArcSettings,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_name",
				Description: "The cluster name of the resource.",
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
				Description: "The provisioning state of the ArcSetting proxy resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ArcSettingProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arc_instance_resource_group",
				Description: "The resource group that hosts the Arc agents, i.e. Hybrid Compute Machine resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ArcSettingProperties.ArcInstanceResourceGroup"),
			},
			{
				Name:        "aggregate_state",
				Description: "Aggregate state of Arc agent across the nodes in this HCI cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ArcSettingProperties.AggregateState"),
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
				Name:        "extensions",
				Description: "List of extensions of the arc setting.",
				Hydrate:     listArcExtensions,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "per_node_details",
				Description: "State of Arc agent in each of the nodes.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ArcSettingProperties.PerNodeDetails"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
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

type arcSettingMap struct {
	azurestackhci.ArcSetting
	ClusterName string
}

//// LIST FUNCTION

func listArcSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(azurestackhci.Cluster)
	resourceGroup := strings.Split(*cluster.ID, "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := azurestackhci.NewArcSettingsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListByCluster(ctx, resourceGroup, *cluster.Name)
	if err != nil {
		plugin.Logger(ctx).Error("listArcSettings", "list", err)
		return nil, err
	}

	for _, arcSetting := range result.Values() {
		d.StreamListItem(ctx, arcSettingMap{arcSetting, *cluster.Name})
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listArcSettings", "list_paging", err)
			return nil, err
		}
		for _, arcSetting := range result.Values() {
			d.StreamListItem(ctx, arcSettingMap{arcSetting, *cluster.Name})
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getArcSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getArcSetting")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	clusterName := d.KeyColumnQuals["cluster_name"].GetStringValue()

	// Handle empty name, resourceGroup or clusterName
	if name == "" || resourceGroup == "" || clusterName == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := azurestackhci.NewArcSettingsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	arcSetting, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		plugin.Logger(ctx).Error("getArcSetting", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if arcSetting.ID != nil {
		return arcSettingMap{arcSetting, clusterName}, nil
	}

	return nil, nil
}

func listArcExtensions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arcSetting := h.Item.(arcSettingMap)
	resourceGroup := strings.Split(*arcSetting.ID, "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := azurestackhci.NewExtensionsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	extensions := []azurestackhci.Extension{}

	result, err := client.ListByArcSetting(ctx, resourceGroup, arcSetting.ClusterName, *arcSetting.Name)
	if err != nil {
		plugin.Logger(ctx).Error("listArcSettings", "list", err)
		return nil, err
	}

	extensions = append(extensions, result.Values()...)

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listArcSettings", "list_paging", err)
			return nil, err
		}
		extensions = append(extensions, result.Values()...)
	}

	return extensions, nil
}
