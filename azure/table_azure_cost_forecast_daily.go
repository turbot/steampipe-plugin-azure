package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			costManagementColumns([]*plugin.Column{
				{
					Name:        "mean_value",
					Description: "The forecasted cost value.",
					Type:        proto.ColumnType_DOUBLE,
					Transform:   transform.FromField("PreTaxCostAmount"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostForecastDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	forecastDef, scope, err := buildForecastQueryInput(ctx, d, "Daily")
	if err != nil {
		return nil, err
	}

	// Get session
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cost_forecast_daily.listCostForecastDaily", "connection_error", err)
		return nil, err
	}

	// Get forecast client
	client, err := armcostmanagement.NewForecastClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cost_forecast_daily.listCostForecastDaily", "client_error", err)
		return nil, err
	}

	// Get forecast
	result, err := client.Usage(ctx, scope, forecastDef, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cost_forecast_daily.listCostForecastDaily", "api_error", err)
		return nil, err
	}

	err = streamForecastResults(ctx, d, &result, scope, "Daily")
	if err != nil {
		return nil, err
	}
	return nil, nil
}
