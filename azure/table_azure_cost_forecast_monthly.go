package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
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

	// Get session
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cost_forecast_monthly.listCostForecastMonthly", "connection_error", err)
		return nil, err
	}

	// Get forecast client
	client, err := armcostmanagement.NewForecastClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cost_forecast_monthly.listCostForecastMonthly", "client_error", err)
		return nil, err
	}

	// Get forecast
	result, err := client.Usage(ctx, scope, forecastDef, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cost_forecast_monthly.listCostForecastMonthly", "api_error", err)
		return nil, err
	}

	err = streamForecastResults(ctx, d, &result, scope, "Monthly")
	if err != nil {
		return nil, err
	}
	return nil, nil
}
