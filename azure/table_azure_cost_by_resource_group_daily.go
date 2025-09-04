package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureCostByResourceGroupDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cost_by_resource_group_daily",
		Description: "Azure Cost Management - Daily cost by resource group",
		List: &plugin.ListConfig{
			KeyColumns: append(costManagementKeyColumns(),
				&plugin.KeyColumn{
					Name:    "resource_group",
					Require: plugin.Optional,
				},
			),
			Hydrate: listCostByResourceGroupDaily,
			Tags:    map[string]string{"service": "Microsoft.CostManagement", "action": "Query"},
		},
		Columns: azureColumns(
			costManagementColumns([]*plugin.Column{
				{
					Name:        "resource_group",
					Description: "The name of the Azure resource group",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByResourceGroupDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	queryDef, scope, err := buildCostQueryInput(ctx, d, "DAILY", []string{"ResourceGroupName"})
	if err != nil {
		return nil, err
	}
	return streamCostAndUsage(ctx, d, queryDef, scope, "ResourceGroupName")
}
