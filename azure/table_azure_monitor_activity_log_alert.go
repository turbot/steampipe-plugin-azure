package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureMonitorActivityLogAlert(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_monitor_activity_log_alert",
		Description: "Azure Monitor Activity Log Alert",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMonitorActivityLogAlert,
			Tags: map[string]string{
				"service": "Microsoft.Insights",
				"action":  "activityLogAlerts/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMonitorActivityLogAlerts,
			Tags: map[string]string{
				"service": "Microsoft.Insights",
				"action":  "activityLogAlerts/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the activity log alert rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The Azure ID (ARM ID) for the activity log alert.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the activity log alert.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Description"),
			},
			{
				Name:        "enabled",
				Description: "The flag that indicates whether the activity log alert is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.Enabled"),
			},
			{
				Name:        "scopes",
				Description: "The list of resource id's that this activity log alert is scoped to.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Scopes"),
			},
			{
				Name:        "condition",
				Description: "The condition that will cause this alert to activate.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Condition"),
			},
			{
				Name:        "actions",
				Description: "The actions that will activate when the condition is met.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Actions"),
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

func listMonitorActivityLogAlerts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_activity_log_alert.listMonitorActivityLogAlerts", "session_error", err)
		return nil, err
	}

	client, err := armmonitor.NewActivityLogAlertsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_activity_log_alert.listMonitorActivityLogAlerts", "client_error", err)
		return nil, err
	}

	pager := client.NewListBySubscriptionIDPager(&armmonitor.ActivityLogAlertsClientListBySubscriptionIDOptions{})
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_monitor_activity_log_alert.listMonitorActivityLogAlerts", "api_error", err)
			return nil, err
		}
		for _, alert := range result.Value {
			d.StreamListItem(ctx, alert)
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

func getMonitorActivityLogAlert(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMonitorActivityLogAlert")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_activity_log_alert.getMonitorActivityLogAlert", "session_error", err)
		return nil, err
	}

	client, err := armmonitor.NewActivityLogAlertsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_activity_log_alert.getMonitorActivityLogAlert", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_activity_log_alert.getMonitorActivityLogAlert", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
