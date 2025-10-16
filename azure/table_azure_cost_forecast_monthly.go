package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCostForecastMonthly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cost_forecast_monthly",
		Description: "Azure Cost Management - Monthly cost forecast",
		List: &plugin.ListConfig{
			KeyColumns: costManagementKeyColumns(),
			Hydrate:    listCostForecastMonthly,
			Tags:       map[string]string{"service": "Microsoft.CostManagement", "action": "Forecast"},
		},
		Columns: azureColumns(
			costManagementColumns([]*plugin.Column{}),
		),
	}
}

//// LIST FUNCTION

func listCostForecastMonthly(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	forecastDef, scope, err := buildForecastQueryInput(ctx, d, "Monthly")
	if err != nil {
		return nil, err
	}

	return streamForecastUsage(ctx, d, forecastDef, scope, "Monthly")
}
