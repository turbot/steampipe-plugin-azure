package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/links"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITON

func tableAzureResourceLink(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_resource_link",
		Description: "Azure Resource Link",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			Hydrate:           getResourceLink,
			ShouldIgnoreError: isNotFoundError([]string{"MissingSubscription", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listResourceLinks,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource link.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The fully qualified ID of the resource link.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The resource link type.",
			},
			{
				Name:        "source_id",
				Type:        proto.ColumnType_STRING,
				Description: "The fully qualified ID of the source resource in the link.",
				Transform:   transform.FromField("Properties.SourceID"),
			},
			{
				Name:        "target_id",
				Type:        proto.ColumnType_STRING,
				Description: "The fully qualified ID of the target resource in the link.",
				Transform:   transform.FromField("Properties.TargetID"),
			},
			{
				Name:        "notes",
				Type:        proto.ColumnType_STRING,
				Description: "Notes about the resource link.",
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listResourceLinks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, nil
	}
	subscriptionID := session.SubscriptionID

	resourceLinkClient := links.NewResourceLinksClient(subscriptionID)
	resourceLinkClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := resourceLinkClient.ListAtSubscription(ctx, "")
		if err != nil {
			return nil, err
		}

		for _, resourceLink := range result.Values() {
			d.StreamListItem(ctx, resourceLink)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getResourceLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getResourceLink")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	linkID := d.KeyColumnQuals["id"].GetStringValue()
	if linkID == "" {
		return nil, nil
	}

	resourceLinkClient := links.NewResourceLinksClient(subscriptionID)
	resourceLinkClient.Authorizer = session.Authorizer

	op, err := resourceLinkClient.Get(ctx, linkID)
	if err != nil {
		return nil, nil
	}

	return op, nil
}
