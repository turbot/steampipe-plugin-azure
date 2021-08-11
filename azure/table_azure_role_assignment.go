package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureIamRoleAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_role_assignment",
		Description: "Azure Role Assignment",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			Hydrate:           getIamRoleAssignment,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listIamRoleAssignments,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the role assignment.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a role assignment uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "scope",
				Description: "Current state of the role assignment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.Scope"),
			},
			{
				Name:        "type",
				Description: "Contains the resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "principal_id",
				Description: "Contains the principal id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.PrincipalID"),
			},
			{
				Name:        "principal_type",
				Description: "Principal type of the assigned principal ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.PrincipalType").Transform(transform.ToString),
			},
			{
				Name:        "role_definition_id",
				Description: "Name of the assigned role definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.RoleDefinitionID"),
			},
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

func listIamRoleAssignments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	authorizationClient := authorization.NewRoleAssignmentsClient(subscriptionID)
	authorizationClient.Authorizer = session.Authorizer
	result, err := authorizationClient.List(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, roleAssignment := range result.Values() {
		d.StreamListItem(ctx, roleAssignment)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, roleAssignment := range result.Values() {
			d.StreamListItem(ctx, roleAssignment)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamRoleAssignment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamRoleAssignment")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	roleAssignmentID := d.KeyColumnQuals["id"].GetStringValue()

	authorizationClient := authorization.NewRoleAssignmentsClient(subscriptionID)
	authorizationClient.Authorizer = session.Authorizer

	op, err := authorizationClient.GetByID(ctx, roleAssignmentID)
	if err != nil {
		return nil, err
	}

	return op, nil
}
