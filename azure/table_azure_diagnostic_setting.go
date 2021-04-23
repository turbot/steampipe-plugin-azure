package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureDiagnosticSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_diagnostic_setting",
		Description: "Azure Diagnostic Setting",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getDiagnosticSetting,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listDiagnosticSettings,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource.",
			},
			{
				Name:        "id",
				Description: "The resource Id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_account_id",
				Description: "The resource ID of the storage account to which you would like to send Diagnostic Logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiagnosticSettings.StorageAccountID"),
			},
			{
				Name:        "service_bus_rule_id",
				Description: "The service bus rule Id of the diagnostic setting.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiagnosticSettings.ServiceBusRuleID"),
			},
			{
				Name:        "event_hub_authorization_rule_id",
				Description: "The resource Id for the event hub authorization rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiagnosticSettings.EventHubAuthorizationRuleID"),
			},
			{
				Name:        "event_hub_name",
				Description: "The name of the event hub.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiagnosticSettings.EventHubName"),
			},
			{
				Name:        "workspace_id",
				Description: "The full ARM resource ID of the Log Analytics workspace to which you would like to send Diagnostic Logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiagnosticSettings.WorkspaceID"),
			},
			{
				Name:        "log_analytics_destination_type",
				Description: "A string indicating whether the export to Log Analytics should use the default destinatio type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiagnosticSettings.LogAnalyticsDestinationType"),
			},
			{
				Name:        "metrics",
				Description: "The list of metric settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DiagnosticSettings.Metrics"),
			},
			{
				Name:        "logs",
				Description: "The list of logs settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DiagnosticSettings.Logs"),
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
				Hydrate:     getDiagnosticSettingResourceGroup,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(diagnosticSettingSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listDiagnosticSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	diagnosticSettingClient := insights.NewDiagnosticSettingsClient(subscriptionID)
	diagnosticSettingClient.Authorizer = session.Authorizer

	resourceURI := "/subscriptions/" + subscriptionID
	result, err := diagnosticSettingClient.List(ctx, resourceURI)
	if err != nil {
		return nil, err
	}

	for _, diagnosticSetting := range *result.Value {
		d.StreamListItem(ctx, diagnosticSetting)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDiagnosticSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDiagnosticSetting")

	name := d.KeyColumnQuals["name"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	diagnosticSettingClient := insights.NewDiagnosticSettingsClient(subscriptionID)
	diagnosticSettingClient.Authorizer = session.Authorizer

	resourceURI := "/subscriptions/" + subscriptionID
	op, err := diagnosticSettingClient.Get(ctx, resourceURI, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func diagnosticSettingSubscriptionID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	subscriptionid := strings.Split(id, "/")[1]
	return subscriptionid, nil
}

func getDiagnosticSettingResourceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDiagnosticSettingResourceGroup")

	var resourseGroupID string
	storage_account_id := h.Item.(insights.DiagnosticSettingsResource).StorageAccountID
	event_hub_authorization_rule_id := h.Item.(insights.DiagnosticSettingsResource).EventHubAuthorizationRuleID
	workspace_id := h.Item.(insights.DiagnosticSettingsResource).WorkspaceID

	if storage_account_id != nil {
		resourseGroupID = strings.Split(*storage_account_id, "/")[4]
	} else if event_hub_authorization_rule_id != nil {
		resourseGroupID = strings.Split(*event_hub_authorization_rule_id, "/")[4]
	} else {
		resourseGroupID = strings.Split(*workspace_id, "/")[4]
	}

	plugin.Logger(ctx).Trace("resourseGroupID", resourseGroupID)

	return resourseGroupID, nil
}
