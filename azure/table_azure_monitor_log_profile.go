package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
)

//// TABLE DEFINITION

func tableAzureMonitorLogProfile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_monitor_log_profile",
		Description: "Azure Monitor Log Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name"}),
			Hydrate:    getMonitorLogProfile,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMonitorLogProfiles,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Azure resource Id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "Azure resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Azure resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The resource location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_account_id",
				Description: "The resource id of the storage account to which you would like to send the Activity Log.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogProfileProperties.StorageAccountID"),
			},
			{
				Name:        "service_bus_rule_id",
				Description: "The service bus rule ID of the service bus namespace in which you would like to have Event Hubs created for streaming the Activity Log.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogProfileProperties.ServiceBusRuleID"),
			},
			{
				Name:        "locations",
				Description: "List of regions for which Activity Log events should be stored or streamed. It is a comma separated list of valid ARM locations including the 'global' location.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LogProfileProperties.Locations"),
			},
			{
				Name:        "categories",
				Description: "The categories of the logs. These categories are created as is convenient to the user.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LogProfileProperties.Categories"),
			},
			{
				Name:        "retention_policy",
				Description: "The retention policy for the events in the log.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LogProfileProperties.RetentionPolicy"),
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

			// Azure standard columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listMonitorLogProfiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_log_profile.listMonitorLogProfiles", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewLogProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// API doesn't support pagination
	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_log_profile.listMonitorLogProfiles", "api_error", err)
		return nil, err
	}

	for _, profile := range *result.Value {
		d.StreamListItem(ctx, profile)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMonitorLogProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_log_profile.getMonitorLogProfile", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewLogProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_monitor_log_profile.getMonitorLogProfile", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
