package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databricks/armdatabricks"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
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
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listDatabricksWorkspaces,
			Tags: map[string]string{
				"service": "Microsoft.Databricks",
				"action":  "workspaces/read",
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listDatabricksWorkspaceDiagnosticSettings,
				Tags: map[string]string{
					"service": "Microsoft.Databricks",
					"action":  "workspaces/providers/Microsoft.Insights/diagnosticSettings/read",
				},
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
				Transform:   transform.FromField("Properties.ManagedResourceGroupID"),
			},
			{
				Name:        "parameters",
				Description: "The workspace's custom parameters.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Parameters"),
			},
			{
				Name:        "provisioning_state",
				Description: "The workspace provisioning state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "ui_definition_uri",
				Description: "The blob URI where the UI definition file is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.UIDefinitionURI"),
			},
			{
				Name:        "authorizations",
				Description: "The workspace provider authorizations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Authorizations"),
			},
			{
				Name:        "created_by",
				Description: "Indicates the Object ID, PUID and Application ID of entity that created the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CreatedBy"),
			},
			{
				Name:        "updated_by",
				Description: "Indicates the Object ID, PUID and Application ID of entity that last updated the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.UpdatedBy"),
			},
			{
				Name:        "created_date_time",
				Description: "Specifies the date and time when the workspace is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.CreatedDateTime"),
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier of the databricks workspace in databricks control plane.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.WorkspaceID"),
			},
			{
				Name:        "workspace_url",
				Description: "The workspace URL which is of the format 'adb-{workspaceId}.{random}.azuredatabricks.net'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.WorkspaceURL"),
			},
			{
				Name:        "storage_account_identity",
				Description: "The details of Managed Identity of Storage Account",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.StorageAccountIdentity"),
			},
			{
				Name:        "public_network_access",
				Description: "The network access type for accessing workspace. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PublicNetworkAccess"),
			},
			{
				Name:        "required_nsg_rules",
				Description: "A value indicating whether the workspace requires compliance with NSG rules. Possible values include: 'AllRules', 'NoAzureDatabricksRules', 'NoAzureServiceRules'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RequiredNsgRules"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of diagnostic settings for the databricks workspace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listDatabricksWorkspaceDiagnosticSettings,
				Transform:   transform.FromValue(),
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
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_databricks_workspace.listDatabricksWorkspaces", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := armdatabricks.NewWorkspacesClient(subscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_databricks_workspace.listDatabricksWorkspaces", "client_error", err)
		return nil, err
	}

	pager := client.NewListBySubscriptionPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_databricks_workspace.listDatabricksWorkspaces", "api_error", err)
			return nil, err
		}

		for _, workspace := range page.Value {
			d.StreamListItem(ctx, workspace)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDatabricksWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("azure_databricks_workspace.getDatabricksWorkspace")

	workspaceName := d.EqualsQualString("name")
	resourceGroup := d.EqualsQualString("resource_group")

	// Return nil, if no input provide
	if workspaceName == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := armdatabricks.NewWorkspacesClient(subscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_databricks_workspace.getDatabricksWorkspace", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, workspaceName, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_databricks_workspace.getDatabricksWorkspace", "api_error", err)
		return nil, err
	}

	return op.Workspace, nil
}

func listDatabricksWorkspaceDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listDatabricksWorkspaceDiagnosticSettings")
	data := h.Item.(*armdatabricks.Workspace)
	id := *data.ID

	// Create session
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armmonitor.NewDiagnosticSettingsClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_databricks_workspace.listDatabricksWorkspaceDiagnosticSettings", "client_error", err)
		return nil, err
	}

	var diagnosticSettings []map[string]interface{}

	input := &armmonitor.DiagnosticSettingsClientListOptions{}

	pager := clientFactory.NewListPager(id, input)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_databricks_workspace.listDatabricksWorkspaceDiagnosticSettings", "api_error", err)
			return nil, err
		}
		for _, i := range page.Value {
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
			if i.Properties != nil {
				objectMap["properties"] = i.Properties
			}
			diagnosticSettings = append(diagnosticSettings, objectMap)
		}
	}
	return diagnosticSettings, nil
}
