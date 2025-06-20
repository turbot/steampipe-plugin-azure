package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureMonitorDiagnosticSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_monitor_diagnostic_setting",
		Description: "Azure Monitor Diagnostic Setting",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_uri"}),
			Hydrate:    getMonitorDiagnosticSetting,
			Tags: map[string]string{
				"service": "Microsoft.Insights",
				"action":  "diagnosticSettings/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMonitorDiagnosticSettings,
			Tags: map[string]string{
				"service": "Microsoft.Insights",
				"action":  "diagnosticSettings/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the diagnostic setting.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID of the diagnostic setting.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_uri",
				Description: "The resource URI for which diagnostic setting needs to be created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_uri"),
			},
			{
				Name:        "storage_account_id",
				Description: "The resource ID of the storage account to which you would like to send diagnostic logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.StorageAccountID"),
			},
			{
				Name:        "service_bus_rule_id",
				Description: "The service bus rule ID of the service bus namespace in which you would like to have Event Hubs created for streaming the diagnostic logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ServiceBusRuleID"),
			},
			{
				Name:        "event_hub_authorization_rule_id",
				Description: "The resource Id for the event hub authorization rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EventHubAuthorizationRuleID"),
			},
			{
				Name:        "event_hub_name",
				Description: "The name of the event hub.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EventHubName"),
			},
			{
				Name:        "metrics",
				Description: "The list of metric settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Metrics"),
			},
			{
				Name:        "logs",
				Description: "The list of logs settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Logs"),
			},
			{
				Name:        "workspace_id",
				Description: "The workspace ID (resource ID of a Log Analytics workspace) for a Log Analytics workspace to which you would like to send diagnostic logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.WorkspaceID"),
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
		}),
	}
}

//// LIST FUNCTION

func listMonitorDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	resourceURI := d.EqualsQuals["resource_uri"].GetStringValue()
	if resourceURI == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_diagnostic_setting.listMonitorDiagnosticSettings", "session_error", err)
		return nil, err
	}

	client, err := armmonitor.NewDiagnosticSettingsClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_diagnostic_setting.listMonitorDiagnosticSettings", "client_error", err)
		return nil, err
	}

	pager := client.NewListPager(resourceURI, nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_monitor_diagnostic_setting.listMonitorDiagnosticSettings", "api_error", err)
			return nil, err
		}
		for _, setting := range result.Value {
			d.StreamListItem(ctx, setting)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMonitorDiagnosticSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMonitorDiagnosticSetting")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceURI := d.EqualsQuals["resource_uri"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || resourceURI == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_diagnostic_setting.getMonitorDiagnosticSetting", "session_error", err)
		return nil, err
	}

	client, err := armmonitor.NewDiagnosticSettingsClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_diagnostic_setting.getMonitorDiagnosticSetting", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceURI, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_diagnostic_setting.getMonitorDiagnosticSetting", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op.DiagnosticSettingsResource, nil
	}

	return nil, nil
}
