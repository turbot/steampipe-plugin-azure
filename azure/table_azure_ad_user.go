package azure

import (
	"context"
	"errors"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureAdUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_ad_user",
		Description: "[DEPRECATED] This table has been deprecated and will be removed in a future release. Please use the azuread_user table in the azuread plugin instead.",
		List: &plugin.ListConfig{
			Hydrate: listAdUsers,
		},

		Columns: []*plugin.Column{
			{
				Name:        "object_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique ID that identifies an active directory user.",
				Transform:   transform.FromField("ObjectID"),
			},
			{
				Name:        "user_principal_name",
				Description: "Principal email of the active directory user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A friendly name that identifies an active directory user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_type",
				Description: "A string that identifies the object type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ObjectType").Transform(transform.ToString),
			},
			{
				Name:        "user_type",
				Description: "A string value that can be used to classify user types in your directory.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserType").Transform(transform.ToString),
			},
			{
				Name:        "given_name",
				Description: "The given name(first name) of the active directory user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "surname",
				Description: "Family name or last name of the active directory user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_enabled",
				Description: "Specifies the account status of the active directory user.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "deletion_timestamp",
				Description: " The time at which the directory object was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "immutable_id",
				Description: "Used to associate an on-premises Active Directory user account with their Azure AD user object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImmutableID"),
			},
			{
				Name:        "mail",
				Description: "The SMTP address for the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail_nickname",
				Description: "The mail alias for the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "usage_location",
				Description: "A two letter country code (ISO standard 3166), required for users that will be assigned licenses due to legal requirement to check for availability of services in countries.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "additional_properties",
				Description: "A list of unmatched properties from the message are deserialized this collection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sign_in_names",
				Description: "A list of sign-in names for a local account in an Azure Active Directory B2C tenant.",
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

func listAdUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	err := errors.New("The azure_ad_user table has been deprecated and removed, please use azuread_user table instead.")
	return nil, err
}
