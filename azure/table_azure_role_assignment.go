package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureIamRoleAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_role_assignment",
		Description: "Azure Role Assignment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIamRoleAssignment,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIamRoleAssignments,
			// For the time being, the optional qualifiers have been commented out
			// due to an issue with the Azure REST API that generates the following error:
			// 	{
			//   "error": {
			//     "code": "UnsupportedQuery",
			//     "message": "The filter 'principalId' is not supported. Supported filters are either 'atScope()' or 'principalId eq '{value}' or assignedTo('{value}')'."
			//   }
			// }
			// Ref: https://github.com/Azure/azure-rest-api-specs/issues/28255
			// We will uncomment it once the issue is resolved.

			// KeyColumns: []*plugin.KeyColumn{
			// 	{
			// 		Name:    "principal_id",
			// 		Require: plugin.Optional,
			// 	},
			// },
		},
		Columns: azureColumns([]*plugin.Column{
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
				Transform:   transform.FromField("Properties.Scope"),
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
				Transform:   transform.FromField("Properties.PrincipalID"),
			},
			{
				Name:        "principal_type",
				Description: "Principal type of the assigned principal ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PrincipalType"),
			},
			{
				Name:        "created_on",
				Description: "Time it was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.CreatedOn"),
			},
			{
				Name:        "updated_on",
				Description: "Time it was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.UpdatedOn"),
			},
			{
				Name:        "role_definition_id",
				Description: "Name of the assigned role definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RoleDefinitionID"),
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
		}),
	}
}

//// LIST FUNCTION

func listIamRoleAssignments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_role_assignment.listIamRoleAssignments", "session_error", err)
		return nil, err
	}

	authorizationClient, err := armauthorization.NewRoleAssignmentsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_role_assignment.listIamRoleAssignments", "client_error", err)
		return nil, err
	}

	defaultFilter := "atScope()" // filter all result
	// Tenant ID is not a required parameter to make the API call.
	option := &armauthorization.RoleAssignmentsClientListForScopeOptions{
		Filter:   &defaultFilter,
	}

	// For the time being, the optional qualifiers have been commented out
	// due to an issue with the Azure REST API that generates the following error:
	// 	{
	//   "error": {
	//     "code": "UnsupportedQuery",
	//     "message": "The filter 'principalId' is not supported. Supported filters are either 'atScope()' or 'principalId eq '{value}' or assignedTo('{value}')'."
	//   }
	// }
	// Ref: https://github.com/Azure/azure-rest-api-specs/issues/28255
	// We will uncomment it once the issue is resolved.

	// var filter string
	// if d.EqualsQuals["principal_id"] != nil {
	// 	filter = fmt.Sprintf("principalId eq '%s'", d.EqualsQuals["principal_id"].GetStringValue())
	// }

	// if filter != "" {
	// 	option.Filter = &filter
	// }

	result := authorizationClient.NewListForScopePager("/subscriptions/"+session.SubscriptionID, option)

	for result.More() {
		res, err := result.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_role_assignment.listIamRoleAssignments", "api_error", err)
			return nil, err
		}
		for _, roleAssignment := range res.Value {
			d.StreamListItem(ctx, roleAssignment)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamRoleAssignment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_role_assignment.getIamRoleAssignment", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	roleAssignmentID := d.EqualsQuals["id"].GetStringValue()

	authorizationClient, err := armauthorization.NewRoleAssignmentsClient(subscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_role_assignment.getIamRoleAssignment", "client_error", err)
		return nil, err
	}

	op, err := authorizationClient.GetByID(ctx, roleAssignmentID, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_role_assignment.getIamRoleAssignment", "api_error", err)
		return nil, err
	}

	return op, nil
}
