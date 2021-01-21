package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAdUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_ad_user",
		Description: "Azure AD User",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("object_id"),
			ItemFromKey:       userObjectIDFromKey,
			Hydrate:           getAdUser,
			ShouldIgnoreError: isNotFoundError([]string{"Request_ResourceNotFound", "Request_BadRequest"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdUsers,
		},

		Columns: []*plugin.Column{
			{
				Name:        "object_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique ID that identifies an active directory user",
				Transform:   transform.FromField("ObjectID"),
			},
			{
				Name:        "user_principal_name",
				Description: "Principal email of the active directory user",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A friendly name that identifies an active directory user",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_type",
				Description: "A string that identifies the object type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ObjectType").Transform(transform.ToString),
			},
			{
				Name:        "user_type",
				Description: "A string value that can be used to classify user types in your directory",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserType").Transform(transform.ToString),
			},
			{
				Name:        "given_name",
				Description: "The given name(first name) of the active directory user",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "surname",
				Description: "Family name or last name of the active directory user",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_enabled",
				Description: "Specifies the account status of the active directory user",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "deletion_timestamp",
				Description: " The time at which the directory object was deleted",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "immutable_id",
				Description: "Used to associate an on-premises Active Directory user account with their Azure AD user object",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImmutableID"),
			},
			{
				Name:        "mail",
				Description: "The SMTP address for the user",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail_nickname",
				Description: "The mail alias for the user",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "usage_location",
				Description: "A two letter country code (ISO standard 3166), required for users that will be assigned licenses due to legal requirement to check for availability of services in countries",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "additional_properties",
				Description: "A list of unmatched properties from the message are deserialized this collection",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sign_in_names",
				Description: "A list of sign-in names for a local account in an Azure Active Directory B2C tenant",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getAdUserTurbotData, "TurbotTitle"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getAdUserTurbotData, "TurbotAkas"),
			},
		},
	}
}

//// ITEM FROM KEY ////

func userObjectIDFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	objectID := quals["object_id"].GetStringValue()
	item := &graphrbac.User{
		ObjectID: &objectID,
	}
	return item, nil
}

//// LIST FUNCTION ////

func listAdUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID

	graphClient := graphrbac.NewUsersClient(tenantID)
	graphClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := graphClient.List(ctx, "", "")
		if err != nil {
			return nil, err
		}

		for _, user := range result.Values() {
			d.StreamListItem(ctx, user)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getAdUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(*graphrbac.User)

	session, err := GetNewSession(ctx, d.ConnectionManager, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID

	graphClient := graphrbac.NewUsersClient(tenantID)
	graphClient.Authorizer = session.Authorizer

	op, err := graphClient.Get(ctx, *user.ObjectID)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS ////

func getAdUserTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(graphrbac.User)
	param := d.Param.(string)

	// Get resource title
	title := data.ObjectID
	if data.DisplayName != nil {
		title = data.DisplayName
	}

	// Get resource tags
	akas := []string{"azure:///user/" + *data.ObjectID}

	if param == "TurbotTitle" {
		return title, nil
	}
	return akas, nil
}
