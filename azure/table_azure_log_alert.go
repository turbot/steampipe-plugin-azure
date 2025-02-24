package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureLogAlert(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_log_alert",
		Description: "Azure Log Alert",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getLogAlert,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLogAlerts,
		},
		Columns: azureColumns([]*plugin.Column{
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
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "Indicates whether this activity log alert is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AlertRuleProperties.Enabled"),
			},
			{
				Name:        "description",
				Description: "A description of this activity log alert.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlertRuleProperties.Description"),
			},
			{
				Name:        "location",
				Description: "The location of the resource. Since Azure Activity Log Alerts is a global service, the location of the rules should always be 'global'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scopes",
				Description: "A list of resourceIds that will be used as prefixes.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AlertRuleProperties.Scopes"),
			},
			{
				Name:        "condition",
				Description: "The condition that will cause this alert to activate.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AlertRuleProperties.Condition"),
			},
			{
				Name:        "actions",
				Description: "The actions that will activate when the condition is met.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AlertRuleProperties.Actions"),
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

func listLogAlerts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	logAlertClient := insights.NewActivityLogAlertsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	logAlertClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &logAlertClient, d.Connection)

	result, err := logAlertClient.ListBySubscriptionID(ctx)
	if err != nil {
		return nil, err
	}

	for _, alertLog := range result.Values() {
		d.StreamListItem(ctx, alertLog)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getLogAlert(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLogAlert")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	logAlertClient := insights.NewActivityLogAlertsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	logAlertClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &logAlertClient, d.Connection)

	op, err := logAlertClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
