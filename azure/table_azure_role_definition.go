package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureIamRoleDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_role_definition",
		Description: "Azure Role Definition",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ItemFromKey:       roleDefinitionNameFromKey,
			Hydrate:           getIamRoleDefinition,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listIamRoleDefinitions,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the role definition",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a role definition uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Contains the resource type",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_name",
				Description: "Current state of the role definition",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleDefinitionProperties.RoleName"),
			},
			{
				Name:        "role_type",
				Description: "Name of the role definition",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleDefinitionProperties.RoleType"),
			},
			{
				Name:        "description",
				Description: "Description of the role definition",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleDefinitionProperties.Description"),
			},
			{
				Name:        "assignable_scopes",
				Description: "A list of assignable scopes for which the role definition can be assigned",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RoleDefinitionProperties.AssignableScopes"),
			},
			{
				Name:        "permissions",
				Description: "A list of actions, which can be accessed",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RoleDefinitionProperties.Permissions"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleDefinitionProperties.RoleName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
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

func roleDefinitionNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &authorization.RoleDefinition{
		Name: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listIamRoleDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	authorizationClient := authorization.NewRoleDefinitionsClient(subscriptionID)
	authorizationClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := authorizationClient.List(ctx, "/subscriptions/"+subscriptionID, "")
		if err != nil {
			return nil, err
		}

		for _, roleDefinition := range result.Values() {
			d.StreamListItem(ctx, roleDefinition)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamRoleDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	roleDefinition := h.Item.(*authorization.RoleDefinition)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	authorizationClient := authorization.NewRoleDefinitionsClient(subscriptionID)
	authorizationClient.Authorizer = session.Authorizer

	op, err := authorizationClient.Get(ctx, "/subscriptions/"+subscriptionID, *roleDefinition.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
