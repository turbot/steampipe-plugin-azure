package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/managedservices/mgmt/managedservices"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureLighthouseDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lighthouse_definition",
		Description: "Azure Lighthouse Definition",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "registration_definition_id",
					Require: plugin.Required,
				},
				{
					Name:    "scope",
					Require: plugin.Optional,
				},
			},
			Hydrate: getAzureLighthouseDefinition,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureLighthouseDefinitions,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"}),
			},
			KeyColumns: plugin.OptionalColumns([]string{"scope"}),
		},
		Columns: []*plugin.Column{
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
				Transform:   transform.FromQual("registration_definition_id"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "scope",
				Description: "The scope of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("scope"),
			},
			{
				Name:        "name",
				Description: "Name of the registration definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
				Name:        "provisioning_state",
				Description: "Current state of the registration definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "authorizations",
				Description: "Authorization details containing principal ID and role ID.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Authorizations"),
			},
			{
				Name:        "plan",
				Description: "Plan details for the managed services.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Plan"),
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		},
	}
}

//// LIST FUNCTION

func listAzureLighthouseDefinitions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.listAzureLighthouseDefinitions", "session_error", err)
		return nil, err
	}
	client := managedservices.NewRegistrationDefinitionsClientWithBaseURI(session.ResourceManagerEndpoint)
	client.Authorizer = session.Authorizer

	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "subscription/" + session.SubscriptionID
	}

	result, err := client.List(ctx, scope)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.listAzureLighthouseDefinitions", "api_error", err)
		return nil, err
	}
	for _, definition := range result.Values() {
		d.StreamListItem(ctx, definition)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_lighthouse_definition.listAzureLighthouseDefinitions", "api_paging_error", err)
			return nil, err
		}
		for _, defn := range result.Values() {
			d.StreamListItem(ctx, defn)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
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

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.getAzureLighthouseDefinition", "session_error", err)
		return nil, err
	}
	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "subscription/" + session.SubscriptionID
	}

	client := managedservices.NewRegistrationDefinitionsClientWithBaseURI(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.Get(ctx, scope, id)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_definition.getAzureLighthouseDefinition", "api_error", err)
		return nil, err
	}
	return result, nil
}
