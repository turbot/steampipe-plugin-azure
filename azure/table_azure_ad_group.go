package azure

import (
	"context"
	"errors"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureAdGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_ad_group",
		Description: "[DEPRECATED] This table has been deprecated and will be removed in a future release. Please use the azuread_group table in the azuread plugin instead.",
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

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listAdGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	err := errors.New("The azure_ad_group table has been deprecated and removed, please use azuread_group table instead.")
	return nil, err
}
