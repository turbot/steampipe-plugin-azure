package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/alertsmanagement/mgmt/alertsmanagement"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAlertMangement(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_alert_management",
		Description: "Azure Alert Management Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getAlertManagement,
			Tags: map[string]string{
				"service": "Microsoft.AlertsManagement",
				"action":  "alerts/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "InvalidApiVersionParameter", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAlertManagements,
			Tags: map[string]string{
				"service": "Microsoft.AlertsManagement",
				"action":  "alerts/read",
			},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "target_resource",
					Require: plugin.Optional,
				},
				{
					Name:    "target_resource_type",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_group",
					Require: plugin.Optional,
				},
				{
					Name:    "alert_rule",
					Require: plugin.Optional,
				},
				{
					Name:    "smart_group_id",
					Require: plugin.Optional,
				},
				{
					Name:    "sort_order",
					Require: plugin.Optional,
				},
				{
					Name:    "custom_time_range",
					Require: plugin.Optional,
				},
				{
					Name:    "sort_by",
					Require: plugin.Optional,
				},
				{
					Name:    "monitor_service",
					Require: plugin.Optional,
				},
				{
					Name:    "monitor_condition",
					Require: plugin.Optional,
				},
				{
					Name:    "severity",
					Require: plugin.Optional,
				},
				{
					Name:    "alert_state",
					Require: plugin.Optional,
				},
				{
					Name:    "time_range",
					Require: plugin.Optional,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "A friendly name that identifies an Alert management service.",
			},
			{
				Name:        "id",
				Description: "Azure resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sort_order",
				Description: "Sort order of the alert management.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("sort_order"),
			},
			{
				Name:        "sort_by",
				Description: "Sort the query results by input field, default value is 'lastModifiedDateTime'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("sort_by"),
			},
			{
				Name:        "custom_time_range",
				Description: "Filter by custom time range in the format <start-time>/<end-time> where time is in (ISO-8601 format).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("custom_time_range"),
			},
			{
				Name:        "time_range",
				Description: "Filter by time range. Possible values are '1h', '1d', '7d' or '30d'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("time_range"),
			},
			{
				Name:        "severity",
				Description: "Severity of alert Sev0 being highest and Sev4 being lowest. Possible values include: 'Sev0', 'Sev1', 'Sev2', 'Sev3', 'Sev4'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.Severity"),
			},
			{
				Name:        "signal_type",
				Description: "The type of signal the alert is based on, which could be metrics, logs or activity logs. Possible values include: 'Metric', 'Log', 'Unknown'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.SignalType"),
			},
			{
				Name:        "alert_state",
				Description: "Alert object state, which can be modified by the user. Possible values include: 'AlertStateNew', 'AlertStateAcknowledged', 'AlertStateClosed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.AlertState"),
			},
			{
				Name:        "monitor_condition",
				Description: "Can be 'Fired' or 'Resolved', which represents whether the underlying conditions have crossed the defined alert rule thresholds. Possible values include: 'Fired', 'Resolved'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.MonitorCondition"),
			},
			{
				Name:        "target_resource",
				Description: "Target ARM resource, on which alert got created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.TargetResource"),
			},
			{
				Name:        "target_resource_name",
				Description: "Name of the target ARM resource, on which alert got created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.TargetResourceName"),
			},
			{
				Name:        "target_resource_type",
				Description: "Resource type of target ARM resource, on which alert got created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.TargetResourceType"),
			},
			{
				Name:        "monitor_service",
				Description: "Monitor service on which the rule(monitor) is set. Possible values include: 'ApplicationInsights', 'ActivityLogAdministrative', 'ActivityLogSecurity', 'ActivityLogRecommendation', 'ActivityLogPolicy', 'ActivityLogAutoscale', 'LogAnalytics', 'Nagios', 'Platform', 'SCOM', 'ServiceHealth', 'SmartDetector', 'VMInsights', 'Zabbix', 'ResourceHealth'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.MonitorService"),
			},
			{
				Name:        "alert_rule",
				Description: "Rule(monitor) which fired alert instance. Depending on the monitor service, this would be ARM ID or name of the rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.AlertRule"),
			},
			{
				Name:        "source_created_id",
				Description: "Unique ID created by monitor service for each alert instance. This could be used to track the issue at the monitor service, in case of Nagios, Zabbix, SCOM, etc.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.SourceCreatedID"),
			},
			{
				Name:        "smart_group_id",
				Description: "Unique ID of the smart group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.SmartGroupID"),
			},
			{
				Name:        "smart_grouping_reason",
				Description: "Verbose reason describing the reason why this alert instance is added to a smart group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.SmartGroupingReason"),
			},
			{
				Name:        "start_date_time",
				Description: "Creation time(ISO-8601 format) of alert instance.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.Essentials.StartDateTime").Transform(convertDateToTime),
			},
			{
				Name:        "last_modified_date_time",
				Description: "Last modification time(ISO-8601 format) of alert instance.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.Essentials.LastModifiedDateTime").Transform(convertDateToTime),
			},
			{
				Name:        "monitor_condition_resolved_date_time",
				Description: "Resolved time(ISO-8601 format) of alert instance. This will be updated when monitor service resolves the alert instance because the rule condition is no longer met.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.Essentials.MonitorConditionResolvedDateTime").Transform(convertDateToTime),
			},
			{
				Name:        "last_modified_user_name",
				Description: "User who last modified the alert, in case of monitor service updates user would be 'system', otherwise name of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Essentials.LastModifiedUserName"),
			},
			{
				Name:        "context",
				Description: "The context of the alert management.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Context"),
			},
			{
				Name:        "egress_config",
				Description: "The egress config for the context management.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EgressConfig"),
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
				Transform:   transform.FromQual("resource_group"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAlertManagements(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	alertManagementClient := alertsmanagement.NewAlertsClientWithBaseURI(session.ResourceManagerEndpoint, "subscriptions/"+subscriptionID, subscriptionID, "")
	alertManagementClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &alertManagementClient, d.Connection)

	var targetResource, targetResourceType, targetResourceGroup, alertRule, smartGroupID, sortOrder, selectParameter, customTimeRange string
	var includeContext, includeEgressConfig bool = true, true
	var pageCount *int32
	var sortBy alertsmanagement.AlertsSortByFields
	var timeRange alertsmanagement.TimeRange
	var monitorService alertsmanagement.MonitorService
	var monitorCondition alertsmanagement.MonitorCondition
	var severity alertsmanagement.Severity
	var alertState alertsmanagement.AlertState

	targetResource = d.EqualsQualString("target_resource")
	targetResourceType = d.EqualsQualString("target_resource_type")
	targetResourceGroup = d.EqualsQualString("resource_group")
	alertRule = d.EqualsQualString("alert_rule")
	smartGroupID = d.EqualsQualString("smart_group_id")
	sortOrder = d.EqualsQualString("sort_order")
	customTimeRange = d.EqualsQualString("custom_time_range")
	sortBy = getShortOrderValue(d.EqualsQualString("sort_by"))
	monitorService = getMonitorServiceValue(d.EqualsQualString("monitor_service"))
	monitorCondition = getMonitorConditionValue(d.EqualsQualString("monitor_condition"))
	severity = getSeverityValue(d.EqualsQualString("severity"))
	alertState = getAlertStateValue(d.EqualsQualString("alert_state"))
	timeRange = getAlertTimeRangeValue(d.EqualsQualString("time_range"))

	result, err := alertManagementClient.GetAll(ctx, targetResource, targetResourceType, targetResourceGroup, monitorService, monitorCondition, severity, alertState, alertRule, smartGroupID, &includeContext, &includeEgressConfig, pageCount, sortBy, sortOrder, selectParameter, timeRange, customTimeRange)
	if err != nil {
		return nil, err
	}

	for _, alertManagement := range result.Values() {
		d.StreamListItem(ctx, alertManagement)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, alertManagement := range result.Values() {
			d.StreamListItem(ctx, alertManagement)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAlertManagement(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	alertId := d.EqualsQuals["id"].GetStringValue()
	if len(alertId) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	alertManagementClient := alertsmanagement.NewAlertsClientWithBaseURI(session.ResourceManagerEndpoint, "", subscriptionID, "")
	alertManagementClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &alertManagementClient, d.Connection)

	op, err := alertManagementClient.GetByID(ctx, alertId)
	if err != nil {
		return nil, err
	}

	return op, nil
}

// // INPUT PARAMETER FUNCTIONS
// We currently lack an SDK-defined function for retrieving the enum value based on the enum string value. To achieve this, explicit manipulation is required.
func getShortOrderValue(s string) alertsmanagement.AlertsSortByFields {
	sortByFields := alertsmanagement.PossibleAlertsSortByFieldsValues()
	for _, i := range sortByFields {
		if s == fmt.Sprint(i) {
			return i
		}
	}
	return ""
}

func getMonitorServiceValue(s string) alertsmanagement.MonitorService {
	svc := alertsmanagement.PossibleMonitorServiceValues()
	for _, i := range svc {
		if s == fmt.Sprint(i) {
			return i
		}
	}
	return ""
}

func getMonitorConditionValue(s string) alertsmanagement.MonitorCondition {
	con := alertsmanagement.PossibleMonitorConditionValues()
	for _, i := range con {
		if s == fmt.Sprint(i) {
			return i
		}
	}
	return ""
}

func getSeverityValue(s string) alertsmanagement.Severity {
	sev := alertsmanagement.PossibleSeverityValues()
	for _, i := range sev {
		if s == fmt.Sprint(i) {
			return i
		}
	}
	return ""
}

func getAlertStateValue(s string) alertsmanagement.AlertState {
	alt := alertsmanagement.PossibleAlertStateValues()
	for _, i := range alt {
		if s == fmt.Sprint(i) {
			return i
		}
	}
	return ""
}

func getAlertTimeRangeValue(s string) alertsmanagement.TimeRange {
	t := alertsmanagement.PossibleTimeRangeValues()
	for _, i := range t {
		if s == fmt.Sprint(i) {
			return i
		}
	}
	return ""
}
