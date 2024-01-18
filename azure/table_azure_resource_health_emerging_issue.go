package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcehealth/armresourcehealth"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureResourceHealthEmergingIssue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_resource_health_emerging_issue",
		Description: "Azure Resource Health Emerging Issue",
		List: &plugin.ListConfig{
			Hydrate: listResourceHealthEmergingIssues,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"}),
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "refresh_timestamp",
				Description: "Timestamp for when last time refreshed for ongoing emerging issue.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.RefreshTimestamp").Transform(convertDateToTime).Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "status_banners",
				Description: "The list of emerging issues of banner type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.StatusBanners"),
			},
			{
				Name:        "status_active_events",
				Description: "The list of emerging issues of active event type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.StatusActiveEvents"),
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

func listResourceHealthEmergingIssues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_resource_health_emerging_issue.listResourceHealthEmergingIssues", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	clientFactory, err := armresourcehealth.NewClientFactory(subscriptionID, session.Creds, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_resource_health_emerging_issue.listResourceHealthEmergingIssues", "NewClientFactory", err)
	}

	pager := clientFactory.NewEmergingIssuesClient().NewListPager(nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_resource_health_emerging_issue.listResourceHealthEmergingIssues", "api_paging_error", err)
		}
		for _, v := range page.Value {
			d.StreamListItem(ctx, v)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if page.NextLink == nil {
			break
		}
	}
	return nil, err
}
