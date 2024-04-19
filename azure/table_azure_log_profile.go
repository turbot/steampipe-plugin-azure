package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureLogProfile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_log_profile",
		Description: "Azure Log Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getLogProfile,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLogProfiles,
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
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "Specifies the name of the region, the resource is created at.",
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
				Name:        "log_event_location",
				Description: "List of regions for which Activity Log events should be stored or streamed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LogProfileProperties.Locations"),
			},
			{
				Name:        "categories",
				Description: "The categories of the logs.",
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

func listLogProfiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	logProfileClient := insights.NewLogProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	logProfileClient.Authorizer = session.Authorizer

	result, err := logProfileClient.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, logProfile := range *result.Value {
		d.StreamListItem(ctx, logProfile)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getLogProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLogProfile")

	name := d.EqualsQuals["name"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	logProfileClient := insights.NewLogProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	logProfileClient.Authorizer = session.Authorizer

	op, err := logProfileClient.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
