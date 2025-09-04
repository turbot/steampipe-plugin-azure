package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureCostByServiceMonthly(_ context.Context) *plugin.Table {
	keyColumns := costManagementKeyColumns()
	keyColumns = append(keyColumns, &plugin.KeyColumn{
		Name:      "service_name",
		Operators: []string{"=", "<>"},
		Require:   plugin.Optional,
	})

	return &plugin.Table{
		Name:        "azure_cost_by_service_monthly",
		Description: "Azure Cost Management - Cost by Service (Monthly)",
		List: &plugin.ListConfig{
			Hydrate:    listCostByServiceMonthly,
			Tags:       map[string]string{"service": "Microsoft.CostManagement", "action": "Query"},
			KeyColumns: keyColumns,
		},
		Columns: azureColumns(
			costManagementColumns([]*plugin.Column{
				{
					Name:        "service_name",
					Description: "The Azure service name.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByServiceMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	queryDef, scope, err := buildCostQueryInput(ctx, d, "MONTHLY", []string{"ServiceName"})
	if err != nil {
		return nil, err
	}
	return streamCostAndUsage(ctx, d, queryDef, scope, "ServiceName")
}
