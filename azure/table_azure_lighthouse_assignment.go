package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managedservices/armmanagedservices"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureLighthouseAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_lighthouse_assignment",
		Description: "Azure Lighthouse Assignment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "registration_assignment_id",
					Require: plugin.Required,
				},
				{
					Name:    "scope",
					Require: plugin.Optional,
				},
			},
			Hydrate: getAzureLighthouseAssignment,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureLighthouseAssignments,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"}),
			},
			KeyColumns: plugin.OptionalColumns([]string{"scope"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the registration assignment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified path of the registration assignment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "registration_assignment_id",
				Description: "The ID of the registration assignment.",
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
				Name:        "registration_definition_id",
				Description: "ID of the associated registration definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RegistrationDefinitionID"),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the registration assignment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
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

func listAzureLighthouseAssignments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_assignment.listAzureLighthouseAssignments", "session_error", err)
		return nil, err
	}
	clientFactory, err := armmanagedservices.NewRegistrationAssignmentsClient(session.Cred, session.ClientOptions)

	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "subscriptions/" + session.SubscriptionID
	}

	pager := clientFactory.NewListPager(scope, &armmanagedservices.RegistrationAssignmentsClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_lighthouse_assignment.listAzureLighthouseAssignments", "api_error", err)
			return nil, err
		}
		for _, assignment := range page.Value {
			d.StreamListItem(ctx, assignment)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTION

func getAzureLighthouseAssignment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("registration_assignment_id")

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_assignment.getAzureLighthouseAssignment", "session_error", err)
		return nil, err
	}
	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "subscriptions/" + session.SubscriptionID
	}

	clientFactory, err := armmanagedservices.NewRegistrationAssignmentsClient(session.Cred, session.ClientOptions)

	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_assignment.getAzureLighthouseAssignment", "client_error", err)
		return nil, err
	}

	result, err := clientFactory.Get(ctx, scope, id, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_lighthouse_assignment.getAzureLighthouseAssignment", "api_error", err)
		return nil, err
	}
	return result, nil
}