package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAdServicePrincipal(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_ad_service_principal",
		Description: "Azure AD Service Principal",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("object_id"),
			Hydrate:           getAdServicePrincipal,
			ShouldIgnoreError: isNotFoundError([]string{"Request_ResourceNotFound", "Request_BadRequest"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdServicePrincipals,
		},
		Columns: []*plugin.Column{
			{
				Name:        "object_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique ID that identifies a service principal.",
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
				Description: "A friendly name that identifies a service principal.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_enabled",
				Description: "Indicates whether or not the service principal account is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "app_role_assignment_required",
				Description: "Specifies whether an AppRoleAssignment to a user or group is required before Azure AD will issue a user or access token to the application.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "deletion_timestamp",
				Description: "The time at which the directory object was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "error_url",
				Description: "An URL provided by the author of the associated application to report errors when using the application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ErrorURL"),
			},
			{
				Name:        "homepage",
				Description: "The URL to the homepage of the associated application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "logout_url",
				Description: "An URL provided by the author of the associated application to logout.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogoutURL"),
			},
			{
				Name:        "saml_metadata_url",
				Description: "The URL to the SAML metadata of the associated application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SamlMetadataURL"),
			},
			{
				Name:        "additional_properties",
				Description: "A list of unmatched properties from the message are deserialized this collection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "alternative_names",
				Description: "A list of alternative names.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "app_roles",
				Description: "A list of application roles that an application may declare. These roles can be assigned to users, groups or service principals.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "key_credentials",
				Description: "A list of key credentials associated with the service principal.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "oauth2_permissions",
				Description: "The OAuth 2.0 permissions exposed by the associated application.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "password_credentials",
				Description: "A list of password credentials associated with the service principal.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "reply_urls",
				Description: "The URLs that user tokens are sent to for sign in with the associated application. The redirect URIs that the oAuth 2.0 authorization code and access tokens are sent to for the associated application.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_principal_names",
				Description: "A list of service principal names.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getAdServicePrincipalTurbotData, "TurbotTitle"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getAdServicePrincipalTurbotData, "TurbotAkas"),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listAdServicePrincipals(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID

	graphClient := graphrbac.NewServicePrincipalsClientWithBaseURI(session.ResourceManagerEndpoint, tenantID)
	graphClient.Authorizer = session.Authorizer

	result, err := graphClient.List(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, servicePrincipal := range result.Values() {
		d.StreamListItem(ctx, servicePrincipal)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, servicePrincipal := range result.Values() {
			d.StreamListItem(ctx, servicePrincipal)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAdServicePrincipal(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAdServicePrincipal")

	session, err := GetNewSession(ctx, d, "GRAPH")
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID
	objectID := d.KeyColumnQuals["object_id"].GetStringValue()

	graphClient := graphrbac.NewServicePrincipalsClientWithBaseURI(session.ResourceManagerEndpoint, tenantID)
	graphClient.Authorizer = session.Authorizer

	op, err := graphClient.Get(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getAdServicePrincipalTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(graphrbac.ServicePrincipal)
	param := d.Param.(string)

	// Get resource title
	title := data.ObjectID
	if data.DisplayName != nil {
		title = data.DisplayName
	}

	// Get resource tags
	akas := []string{"azure:///serviceprincipal/" + *data.ObjectID}

	if param == "TurbotTitle" {
		return title, nil
	}
	return akas, nil
}
