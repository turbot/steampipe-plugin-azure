package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managedservices/armmanagedservices"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureLighthouseDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lighthouse_definition",
		Description: "Azure Lighthouse Definition",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: isNotFoundError([]string{"SubscriptionNotFound"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "registration_definition_id",
					Require: plugin.Required,
				},
				{
					Name:      "scope",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
			},
			Hydrate: getAzureLighthouseDefinition,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"RegistrationDefinitionNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate:    listAzureLighthouseDefinitions,
			KeyColumns: plugin.OptionalColumns([]string{"scope"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the registration definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified path of the registration definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "registration_definition_id",
				Description: "The ID of the registration definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(lastPathElement),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "The scope of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("scope"),
			},
			{
				Name:        "description",
				Description: "Description of the registration definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Description"),
			},
			{
				Name:        "registration_definition_name",
				Description: "Name of the registration definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RegistrationDefinitionName"),
			},
			{
				Name:        "managed_by_tenant_id",
				Description: "ID of the managedBy tenant.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ManagedByTenantID"),
			},
			{
				Name:        "managed_by_tenant_name",
				Description: "The name of the managedBy tenant.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ManagedByTenantName"),
			},
			{
				Name:        "managed_tenant_name",
				Description: "The name of the managed tenant.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ManageeTenantName"),
			},
			{
				Name:        "authorizations",
				Description: "Authorization details containing principal ID and role ID.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Authorizations"),
			},
			{
				Name:        "eligible_authorizations",
				Description: "The collection of eligible authorization objects describing the just-in-time access Azure Active Directory principals in the managedBy tenant will receive on the delegated resource in the managed tenant.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.EligibleAuthorizations"),
			},
			{
				Name:        "plan",
				Description: "Plan details for the managed services.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
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

			// Azure standard columns
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getLighthouseDefinitionResourceGroup,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listAzureLighthouseDefinitions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.listAzureLighthouseDefinitions", "session_error", err)
		return nil, err
	}
	clientFactory, err := armmanagedservices.NewRegistrationDefinitionsClient(session.Cred, session.ClientOptions)

	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "subscriptions/" + session.SubscriptionID
	}

	pager := clientFactory.NewListPager(scope, &armmanagedservices.RegistrationDefinitionsClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_lighthouse_definition.listAzureLighthouseDefinitions", "api_error", err)
			return nil, err
		}
		for _, definition := range page.Value {
			d.StreamListItem(ctx, definition)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTION

func getAzureLighthouseDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("registration_definition_id")

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.getAzureLighthouseDefinition", "session_error", err)
		return nil, err
	}
	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "subscriptions/" + session.SubscriptionID
	}

	clientFactory, err := armmanagedservices.NewRegistrationDefinitionsClient(session.Cred, session.ClientOptions)

	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.getAzureLighthouseDefinition", "client_error", err)
		return nil, err
	}

	result, err := clientFactory.Get(ctx, scope, id, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.getAzureLighthouseDefinition", "api_error", err)
		return nil, err
	}
	return result, nil
}

// We can have Definition/Assignments in different scopes:
// Subscription: /subscriptions/{subscription-id}
// Resource Groups: /subscriptions/{subscription-id}/resourceGroups/{resource-group-name}
// Management Groups: /providers/Microsoft.Management/managementGroups/{management-group-id}
// Individual Resources: /subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/{resource-provider}/{resource-type}/{resource-name}
func getLighthouseDefinitionResourceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if h.Item == nil {
		return nil, nil
	}

	id := ""
	switch item := h.Item.(type) {
	case *armmanagedservices.RegistrationDefinition:
		id = *item.ID
	case armmanagedservices.RegistrationDefinitionsClientGetResponse:
		id = *item.ID
	}

	if id != "" && strings.Contains(strings.ToLower(id), "/resourcegroups/") {
		return strings.Split(id, "/")[4], nil
	}

	return nil, nil
}
