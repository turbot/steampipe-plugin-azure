package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureIamRoleAssignmentForResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_role_assignment_for_resource",
		Description: "Azure Role Assignment for Resource",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "scope"}),
			ItemFromKey:       roleAssignmentDetailsFromKey,
			Hydrate:           getIamRoleAssignmentForResource,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listIamRoleAssignmentForResources,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the role assignment",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a role assignment uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "scope",
				Description: "Current state of the role assignment",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.Scope"),
			},
			{
				Name:        "type",
				Description: "Contains the resource type",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "can_delegate",
				Description: "Delegation flag for the role assignment",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.CanDelegate"),
				Default:     false,
			},
			{
				Name:        "principal_id",
				Description: "Contains the principal id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.PrincipalID"),
			},
			{
				Name:        "principal_type",
				Description: "The principal type of the assigned principal ID",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.PrincipalType").Transform(transform.ToString),
			},
			{
				Name:        "role_definition_id",
				Description: "Name of the assigned role definition",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleAssignmentPropertiesWithScope.RoleDefinitionID"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getIamRoleAssignmentForResourceTurbotData, "TurbotTitle"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getIamRoleAssignmentForResourceTurbotData, "TurbotAkas"),
			},
			{
				Name:        "subscription_id",
				Description: "The Azure Subscription ID in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// ITEM FROM KEY

func roleAssignmentDetailsFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	roleAssignmentName := quals["name"].GetStringValue()
	scope := quals["scope"].GetStringValue()
	item := &authorization.RoleAssignment{
		Name: &roleAssignmentName,
		RoleAssignmentPropertiesWithScope: &authorization.RoleAssignmentPropertiesWithScope{
			Scope: &scope,
		},
	}
	return item, nil
}

//// LIST FUNCTION

func listIamRoleAssignmentForResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	authorizationClient := authorization.NewRoleAssignmentsClient(subscriptionID)
	authorizationClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := authorizationClient.List(ctx, "")
		if err != nil {
			return nil, err
		}

		for _, roleAssignment := range result.Values() {
			d.StreamListItem(ctx, roleAssignment)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamRoleAssignmentForResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	roleAssignment := h.Item.(*authorization.RoleAssignment)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	authorizationClient := authorization.NewRoleAssignmentsClient(subscriptionID)
	authorizationClient.Authorizer = session.Authorizer

	op, err := authorizationClient.Get(ctx, *roleAssignment.RoleAssignmentPropertiesWithScope.Scope, *roleAssignment.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getIamRoleAssignmentForResourceTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(authorization.RoleAssignment)
	param := d.Param.(string)
	splittedID := strings.Split(string(*data.RoleAssignmentPropertiesWithScope.Scope), "/")

	// Get resource title
	title := splittedID[len(splittedID)-2]

	// Get resource tags
	akas := []string{"azure://" + *data.RoleAssignmentPropertiesWithScope.Scope, "azure://" + strings.ToLower(*data.RoleAssignmentPropertiesWithScope.Scope)}

	if param == "TurbotTitle" {
		return title, nil
	}
	return akas, nil
}
