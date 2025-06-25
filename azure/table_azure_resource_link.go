package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/links"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITON

func tableAzureResourceLink(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_resource_link",
		Description: "Azure Resource Link",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getResourceLink,
			Tags: map[string]string{
				"service": "Microsoft.Resources",
				"action":  "links/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listResourceLinks,
			Tags: map[string]string{
				"service": "Microsoft.Resources",
				"action":  "links/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
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
		}),
	}
}

//// LIST FUNCTION

func listResourceLinks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, nil
	}
	subscriptionID := session.SubscriptionID

	resourceLinkClient := links.NewResourceLinksClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	resourceLinkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &resourceLinkClient, d.Connection)

	result, err := resourceLinkClient.ListAtSubscription(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, resourceLink := range result.Values() {
		d.StreamListItem(ctx, resourceLink)
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
		for _, resourceLink := range result.Values() {
			d.StreamListItem(ctx, resourceLink)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
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
	linkID := d.EqualsQuals["id"].GetStringValue()
	if linkID == "" {
		return nil, nil
	}

	resourceLinkClient := links.NewResourceLinksClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	resourceLinkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &resourceLinkClient, d.Connection)

	op, err := resourceLinkClient.Get(ctx, linkID)
	if err != nil {
		return nil, nil
	}

	return op, nil
}
