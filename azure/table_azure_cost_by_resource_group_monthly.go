package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureCostByResourceGroupMonthly(_ context.Context) *plugin.Table {
	keyColumns := costManagementKeyColumns()
	keyColumns = append(keyColumns, &plugin.KeyColumn{
		Name:      "resource_group",
		Operators: []string{"=", "<>"},
		Require:   plugin.Optional,
	})

	return &plugin.Table{
		Name:        "azure_cost_by_resource_group_monthly",
		Description: "Azure Cost Management - Cost by Resource Group (Monthly)",
		List: &plugin.ListConfig{
			Hydrate:    listCostByResourceGroupMonthly,
			Tags:       map[string]string{"service": "Microsoft.CostManagement", "action": "Query"},
			KeyColumns: keyColumns,
		},
		Columns: azureColumns(
			costManagementColumns([]*plugin.Column{
				{
					Name:        "resource_group",
					Description: "The resource group name.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByResourceGroupMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	queryDef, scope, err := buildCostQueryInput(ctx, d, "MONTHLY", []string{"ResourceGroupName"})
	if err != nil {
		return nil, err
	}
	return streamCostAndUsage(ctx, d, queryDef, scope, "ResourceGroupName")
}
