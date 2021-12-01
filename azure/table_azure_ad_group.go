package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAdGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_ad_group",
		Description: "[DEPRECATED] This table has been deprecated and will be removed in a future release. Please use the azuread_group table in the azuread plugin instead.",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("object_id"),
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
				Description: "The unique ID that identifies a group.",
				Transform:   transform.FromField("ObjectID"),
			},
			{
				Name:        "object_type",
				Description: "A string that identifies the object type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ObjectType").Transform(transform.ToString),
			},
			{
				Name:        "display_name",
				Description: "A friendly name that identifies a group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail",
				Description: "The primary email address of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail_enabled",
				Description: "Indicates whether the group is mail-enabled. Must be false. This is because only pure security groups can be created using the Graph API.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "mail_nickname",
				Description: "The mail alias for the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_timestamp",
				Description: "The time at which the directory object was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "security_enabled",
				Description: "Specifies whether the group is a security group.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "additional_properties",
				Description: "A list of unmatched properties from the message are deserialized this collection.",
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

//// LIST FUNCTION

func listAdGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID

	graphClient := graphrbac.NewGroupsClientWithBaseURI(session.ResourceManagerEndpoint, tenantID)
	graphClient.Authorizer = session.Authorizer

	result, err := graphClient.List(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, group := range result.Values() {
		d.StreamListItem(ctx, group)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, group := range result.Values() {
			d.StreamListItem(ctx, group)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAdGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAdGroup")

	session, err := GetNewSession(ctx, d, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID
	objectID := d.KeyColumnQuals["object_id"].GetStringValue()

	graphClient := graphrbac.NewGroupsClientWithBaseURI(session.ResourceManagerEndpoint, tenantID)
	graphClient.Authorizer = session.Authorizer

	op, err := graphClient.Get(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

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
