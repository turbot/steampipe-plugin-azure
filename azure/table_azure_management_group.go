package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/managementgroups"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureManagementGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_management_group",
		Description: "Azure Management Group.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getManagementGroup,
			Tags: map[string]string{
				"service": "Microsoft.Management",
				"action":  "managementGroups/read",
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listManagementGroups,
			Tags: map[string]string{
				"service": "Microsoft.Management",
				"action":  "managementGroups/read",
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The fully qualified ID for the management group.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "The name of the management group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the management group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The friendly name of the management group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InfoProperties.DisplayName", "Properties.DisplayName"),
			},
			{
				Name:        "tenant_id",
				Description: "The AAD Tenant ID associated with the management group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InfoProperties.TenantID", "Properties.TenantID"),
			},
			{
				Name:        "updated_by",
				Description: "The identity of the principal or process that updated the management group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getManagementGroup,
				Transform:   transform.FromField("Properties.Details.UpdatedBy"),
			},
			{
				Name:        "updated_time",
				Description: "The date and time when this management group was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getManagementGroup,
				Transform:   transform.FromField("Properties.Details.UpdatedTime.Time"),
			},
			{
				Name:        "version",
				Description: "The version number of the management group.",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getManagementGroup,
				Transform:   transform.FromField("Properties.Details.Version"),
			},
			{
				Name:        "children",
				Description: "The list of children of the management group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getManagementGroup,
				Transform:   transform.FromField("Properties.Children"),
			},
			{
				Name:        "parent",
				Description: "The associated parent management group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getManagementGroup,
				Transform:   transform.FromField("Properties.Details.Parent"),
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
		},
	}
}

//// LIST FUNCTION

func listManagementGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	mgClient := managementgroups.NewClient()
	mgClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &mgClient, d.Connection)

	result, err := mgClient.List(ctx, "", "")
	if err != nil {
		plugin.Logger(ctx).Error("listManagementGroups", "list", err)
		return nil, err
	}
	for _, mg := range result.Values() {
		d.StreamListItem(ctx, mg)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, mg := range result.Values() {
			d.StreamListItem(ctx, mg)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getManagementGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getManagementGroup")

	var name string
	if h.Item != nil {
		name = *h.Item.(managementgroups.Info).Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	mgClient := managementgroups.NewClient()
	mgClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &mgClient, d.Connection)

	op, err := mgClient.Get(ctx, name, "children", nil, "", "")
	if err != nil {
		plugin.Logger(ctx).Error("getManagementGroup", "get", err)
		return nil, err
	}

	return op, nil
}
