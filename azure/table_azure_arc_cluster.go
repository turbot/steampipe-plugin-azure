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
		Name:        "azure_arc_cluster",
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
				Name:        "status",
				Description: "Status of the cluster agent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.Status"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.ProvisioningState"),
			},
			{
				Name:        "aad_client_id",
				Description: "App id of cluster AAD identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.AadClientID"),
			},
			{
				Name:        "aad_tenant_id",
				Description: "Tenant id of cluster AAD identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.AadTenantID"),
			},
			{
				Name:        "billing_model",
				Description: "Type of billing applied to the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.BillingModel"),
			},
			{
				Name:        "cloud_id",
				Description: "Immutable resource id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.CloudID"),
			},
			{
				Name:        "cloud_management_endpoint",
				Description: "Endpoint configured for management from the Azure portal.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterProperties.CloudManagementEndpoint"),
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
				Name:        "last_billing_timestamp",
				Description: "Most recent billing meter timestamp.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ClusterProperties.LastBillingTimestamp").Transform(convertDateToTime),
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
				Name:        "last_sync_timestamp",
				Description: "Most recent cluster sync timestamp.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ClusterProperties.LastSyncTimestamp").Transform(convertDateToTime),
			},
			{
				Name:        "location",
				Description: "Location of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "registration_timestamp",
				Description: "First cluster sync timestamp.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ClusterProperties.RegistrationTimestamp").Transform(convertDateToTime),
			},
			{
				Name:        "trial_days_remaining",
				Description: "Number of days remaining in the trial period.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ClusterProperties.TrialDaysRemaining"),
			},
			{
				Name:        "reported_properties",
				Description: "Properties reported by cluster agent..",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterProperties.ReportedProperties"),
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
