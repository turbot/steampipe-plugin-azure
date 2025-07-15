package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/operationalinsights/mgmt/operationalinsights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureLogAnalyticsWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_log_analytics_workspace",
		Description: "Azure Log Analytics Workspace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getLogAnalyticsWorkspace,
			Tags: map[string]string{
				"service": "Microsoft.OperationalInsights",
				"action":  "workspaces/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLogAnalyticsWorkspaces,
			Tags: map[string]string{
				"service": "Microsoft.OperationalInsights",
				"action":  "workspaces/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the Log Analytics workspace.",
			},
			{
				Name:        "id",
				Description: "Contains the unique ID to identify the Log Analytics workspace.",
				Transform:   transform.FromGo(),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the Log Analytics workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the Log Analytics workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku",
				Description: "The SKU (pricing level) of the Log Analytics workspace.",
				Transform:   transform.FromField("WorkspaceProperties.Sku"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "retention_in_days",
				Description: "The retention period for the Log Analytics workspace data in days.",
				Transform:   transform.FromField("WorkspaceProperties.RetentionInDays"),
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the Log Analytics workspace.",
				Transform:   transform.FromField("WorkspaceProperties.ProvisioningState"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workspace_capping",
				Description: "The workspace capping properties.",
				Transform:   transform.FromField("WorkspaceProperties.WorkspaceCapping"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_date",
				Description: "Workspace creation date.",
				Transform:   transform.FromField("WorkspaceProperties.CreatedDate").Transform(transform.NullIfZeroValue),
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "modified_date",
				Description: "Workspace modification date.",
				Transform:   transform.FromField("WorkspaceProperties.ModifiedDate").Transform(transform.NullIfZeroValue),
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "customer_id",
				Description: "Represents the ID associated with the workspace.",
				Transform:   transform.FromField("WorkspaceProperties.CustomerID"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_network_access_for_ingestion",
				Description: "The network access type for accessing Log Analytics ingestion.",
				Transform:   transform.FromField("WorkspaceProperties.PublicNetworkAccessForIngestion"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_network_access_for_query",
				Description: "The network access type for accessing Log Analytics query.",
				Transform:   transform.FromField("WorkspaceProperties.PublicNetworkAccessForQuery"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "force_cmk_for_query",
				Description: "Indicates whether customer managed storage is mandatory for query management.",
				Transform:   transform.FromField("WorkspaceProperties.ForceCmkForQuery"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "private_link_scoped_resources",
				Description: "List of linked private link scope resources.",
				Transform:   transform.FromField("WorkspaceProperties.PrivateLinkScopedResources"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "enable_data_export",
				Description: "Flag that indicates if data should be exported.",
				Transform:   transform.FromField("WorkspaceProperties.Features.EnableDataExport"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "immediate_purge_data_on_30_days",
				Description: "Flag that describes if we want to remove the data after 30 days.",
				Transform:   transform.FromField("WorkspaceProperties.Features.ImmediatePurgeDataOn30Days"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "enable_log_access_using_only_resource_permissions",
				Description: "Flag that indicates which permission to use - resource or workspace or both.",
				Transform:   transform.FromField("WorkspaceProperties.Features.EnableLogAccessUsingOnlyResourcePermissions"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cluster_resource_id",
				Description: "Dedicated LA cluster resourceId that is linked to the workspaces.",
				Transform:   transform.FromField("WorkspaceProperties.Features.ClusterResourceID"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disable_local_auth",
				Description: "Disable Non-AAD based Auth.",
				Transform:   transform.FromField("WorkspaceProperties.Features.DisableLocalAuth"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "tags",
				Description: "The tags assigned to the Log Analytics workspace.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "The title of the Log Analytics workspace.",
				Transform:   transform.FromField("Name"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: "The list of globally unique identifier strings (also known as) for the resource.",
				Transform:   transform.FromField("ID").Transform(idToAkas),
				Type:        proto.ColumnType_JSON,
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: "The region of the Log Analytics workspace.",
				Transform:   transform.FromField("Location").Transform(toLower),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_group",
				Description: "The resource group of the Log Analytics workspace.",
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listLogAnalyticsWorkspaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_log_analytics_workspace.listLogAnalyticsWorkspaces", "connection_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := operationalinsights.NewWorkspacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	result, err := client.List(ctx)
	if err != nil {
		logger.Error("azure_log_analytics_workspace.listLogAnalyticsWorkspaces", "api_error", err)
		return nil, err
	}

	for _, workspace := range *result.Value {
		d.StreamListItem(ctx, workspace)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getLogAnalyticsWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_log_analytics_workspace.getLogAnalyticsWorkspace", "connection_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := operationalinsights.NewWorkspacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		logger.Error("azure_log_analytics_workspace.getLogAnalyticsWorkspace", "api_error", err)
		return nil, err
	}

	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
