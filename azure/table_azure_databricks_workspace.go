package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/databricks/mgmt/databricks"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureDatabricksWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_databricks_workspace",
		Description: "Azure Databricks Workspace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getDatabricksWorkspace,
			Tags: map[string]string{
				"service": "Microsoft.Databricks",
				"action":  "workspaces/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listDatabricksWorkspaces,
			Tags: map[string]string{
				"service": "Microsoft.Databricks",
				"action":  "workspaces/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "sku",
				Description: "The SKU of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The geo-location where the resource lives.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_resource_group_id",
				Description: "The managed resource group ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ManagedResourceGroupID"),
			},
			{
				Name:        "parameters",
				Description: "The workspace's custom parameters.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.Parameters"),
			},
			{
				Name:        "provisioning_state",
				Description: "The workspace provisioning state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ProvisioningState"),
			},
			{
				Name:        "ui_definition_uri",
				Description: "The blob URI where the UI definition file is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.UIDefinitionURI"),
			},
			{
				Name:        "authorizations",
				Description: "The workspace provider authorizations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.Authorizations"),
			},
			{
				Name:        "created_by",
				Description: "Indicates the Object ID, PUID and Application ID of entity that created the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.CreatedBy"),
			},
			{
				Name:        "updated_by",
				Description: "Indicates the Object ID, PUID and Application ID of entity that last updated the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.UpdatedBy"),
			},
			{
				Name:        "created_date_time",
				Description: "Specifies the date and time when the workspace is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("WorkspaceProperties.CreatedDateTime").Transform(convertDateToTime),
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier of the databricks workspace in databricks control plane.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.WorkspaceID"),
			},
			{
				Name:        "workspace_url",
				Description: "The workspace URL which is of the format 'adb-{workspaceId}.{random}.azuredatabricks.net'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.WorkspaceURL"),
			},
			{
				Name:        "storage_account_identity",
				Description: "The details of Managed Identity of Storage Account",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.StorageAccountIdentity"),
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
		}),
	}
}

//// LIST FUNCTION

func listDatabricksWorkspaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listDatabricksWorkspaces")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	workspaceClient := databricks.NewWorkspacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	workspaceClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &workspaceClient, d.Connection)

	result, err := workspaceClient.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_databricks_workspace.listDatabricksWorkspaces", "ListBySubscription", err)
		return nil, err
	}
	for _, workspace := range result.Values() {
		d.StreamListItem(ctx, workspace)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_databricks_workspace.listDatabricksWorkspaces", "ListBySubscription_pagination", err)
			return nil, err
		}
		for _, device := range result.Values() {
			d.StreamListItem(ctx, device)
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDatabricksWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getDatabricksWorkspace")

	workspaceName := d.EqualsQualString("name")
	resourceGroup := d.EqualsQualString("resource_group")

	// Return nil, if no input provide
	if workspaceName == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	workspaceClient := databricks.NewWorkspacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	workspaceClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &workspaceClient, d.Connection)

	op, err := workspaceClient.Get(ctx, resourceGroup, workspaceName)
	if err != nil {
		return nil, err
	}

	return op, nil
}
