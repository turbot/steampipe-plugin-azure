package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureMonitorMetricAlert(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_monitor_metric_alert",
		Description: "Azure Monitor Metric Alert",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMonitorMetricAlert,
			Tags: map[string]string{
				"service": "Microsoft.Insights",
				"action":  "metricAlerts/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMonitorMetricAlerts,
			Tags: map[string]string{
				"service": "Microsoft.Insights",
				"action":  "metricAlerts/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the metric alert rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The Azure ID (ARM ID) for the metric alert.",
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
				Description: "The description of the metric alert.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Description"),
			},
			{
				Name:        "severity",
				Description: "The severity of the metric alert.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Severity"),
			},
			{
				Name:        "enabled",
				Description: "The flag that indicates whether the metric alert is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.Enabled"),
			},
			{
				Name:        "scopes",
				Description: "The list of resource id's that this metric alert is scoped to.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Scopes"),
			},
			{
				Name:        "evaluation_frequency",
				Description: "How often the metric alert is evaluated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EvaluationFrequency"),
			},
			{
				Name:        "window_size",
				Description: "The period of time that is used to monitor alert activity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.WindowSize"),
			},
			{
				Name:        "target_resource_type",
				Description: "The resource type of the target resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.TargetResourceType"),
			},
			{
				Name:        "target_resource_region",
				Description: "The region of the target resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.TargetResourceRegion"),
			},
			{
				Name:        "criteria",
				Description: "The rule criteria that defines the conditions of the alert rule.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Criteria"),
			},
			{
				Name:        "auto_mitigate",
				Description: "The flag that indicates whether the alert should be auto resolved or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.AutoMitigate"),
			},
			{
				Name:        "actions",
				Description: "The actions that will activate when the alert rule's conditions are met.",
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

func listMonitorMetricAlerts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_metric_alert.listMonitorMetricAlerts", "session_error", err)
		return nil, err
	}

	client, err := armmonitor.NewMetricAlertsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_metric_alert.listMonitorMetricAlerts", "client_error", err)
		return nil, err
	}

	pager := client.NewListBySubscriptionPager(&armmonitor.MetricAlertsClientListBySubscriptionOptions{})
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_monitor_metric_alert.listMonitorMetricAlerts", "api_error", err)
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

func getMonitorMetricAlert(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMonitorMetricAlert")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_metric_alert.getMonitorMetricAlert", "session_error", err)
		return nil, err
	}

	client, err := armmonitor.NewMetricAlertsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_metric_alert.getMonitorMetricAlert", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_metric_alert.getMonitorMetricAlert", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
