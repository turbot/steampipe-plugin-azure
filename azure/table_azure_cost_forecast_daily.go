package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCostForecastDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cost_forecast_daily",
		Description: "Azure Cost Management - Daily cost forecast",
		List: &plugin.ListConfig{
			KeyColumns: costManagementKeyColumns(),
			Hydrate:    listCostForecastDaily,
			Tags:       map[string]string{"service": "Microsoft.CostManagement", "action": "Forecast"},
		},
		Columns: azureColumns(
			costManagementColumns([]*plugin.Column{}),
		),
	}
}

//// LIST FUNCTION

func listCostForecastDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	forecastDef, scope, err := buildForecastQueryInput(ctx, d, "Daily")
	if err != nil {
		return nil, err
	}

	return streamForecastUsage(ctx, d, forecastDef, scope, "Daily")
}
