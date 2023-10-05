package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
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
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource.",
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
				Transform:   transform.FromField("EmergingIssue.RefreshTimestamp").Transform(convertDateToTime),
			},
			{
				Name:        "status_banners",
				Description: "The list of emerging issues of banner type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EmergingIssue.StatusBanners"),
			},
			{
				Name:        "status_active_events",
				Description: "The list of emerging issues of active event type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EmergingIssue.StatusActiveEvents"),
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID).Transform(toLower),
			},
		}),
	}
}

//// LIST FUNCTION

func listResourceHealthEmergingIssues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	emergingClient := resourcehealth.NewEmergingIssuesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	emergingClient.Authorizer = session.Authorizer
	result, err := emergingClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range result.Values() {
		d.StreamListItem(ctx, item)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range result.Values() {
			d.StreamListItem(ctx, item)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}
