package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAdGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_ad_group",
		Description: "Azure AD Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("object_id"),
			ItemFromKey:       groupObjectIDFromKey,
			ShouldIgnoreError: isNotFoundError([]string{"Request_ResourceNotFound", "Request_BadRequest"}),
			Hydrate:           getAdGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listAdGroups,
		},

		Columns: []*plugin.Column{
			{
				Name:        "object_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique ID that identifies a group",
				Transform:   transform.FromField("ObjectID"),
			},
			{
				Name:        "object_type",
				Description: "A string that identifies the object type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ObjectType").Transform(transform.ToString),
			},
			{
				Name:        "display_name",
				Description: "A friendly name that identifies a group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail",
				Description: "The primary email address of the group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail_enabled",
				Description: "Indicates whether the group is mail-enabled. Must be false. This is because only pure security groups can be created using the Graph API",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "mail_nickname",
				Description: "The mail alias for the group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_timestamp",
				Description: "The time at which the directory object was deleted",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "security_enabled",
				Description: "Specifies whether the group is a security group",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "additional_properties",
				Description: "A list of unmatched properties from the message are deserialized this collection",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getAdGroupTurbotData, "TurbotTitle"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getAdGroupTurbotData, "TurbotAkas"),
			},
		},
	}
}

//// ITEM FROM KEY ////

func groupObjectIDFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	objectID := quals["object_id"].GetStringValue()
	item := &graphrbac.ADGroup{
		ObjectID: &objectID,
	}
	return item, nil
}

//// LIST FUNCTION ////

func listAdGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID

	graphClient := graphrbac.NewGroupsClient(tenantID)
	graphClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := graphClient.List(ctx, "")
		if err != nil {
			return nil, err
		}

		for _, group := range result.Values() {
			d.StreamListItem(ctx, group)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getAdGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(*graphrbac.ADGroup)

	session, err := GetNewSession(ctx, d, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID

	graphClient := graphrbac.NewGroupsClient(tenantID)
	graphClient.Authorizer = session.Authorizer

	op, err := graphClient.Get(ctx, *group.ObjectID)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS ////

func getAdGroupTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(graphrbac.ADGroup)
	param := d.Param.(string)

	// Get resource title
	title := data.ObjectID
	if data.DisplayName != nil {
		title = data.DisplayName
	}

	// Get resource tags
	akas := []string{"azure:///group/" + *data.ObjectID}

	if param == "TurbotTitle" {
		return title, nil
	}
	return akas, nil
}
