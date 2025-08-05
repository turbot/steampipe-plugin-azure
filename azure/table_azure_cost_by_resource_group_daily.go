package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
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
	params, err := buildCostByResourceGroupDailyInput(ctx, d)
	if err != nil {
		return nil, err
	}
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByResourceGroupDailyInput(ctx context.Context, d *plugin.QueryData) (*AzureCostQueryInput, error) {
	// Get subscription ID (will be handled in streamCostAndUsage if empty)
	subscriptionID := d.EqualsQualString("subscription_id")
	if subscriptionID == "" {
		subscriptionID = "placeholder"
	}

	// Set timeframe and time period using new usage_date logic
	timeframe := armcostmanagement.TimeframeTypeCustom
	timePeriod := &armcostmanagement.QueryTimePeriod{}

	// Get time range from period_start/period_end quals
	startTime, endTime := getPeriodTimeRange(d, "DAILY")

	// Set default time range if no quals provided
	if startTime == "" || endTime == "" {
		defaultEnd := time.Now()
		defaultStart := defaultEnd.AddDate(0, -11, -30) // Last 1 year
		if startTime == "" {
			startTime = defaultStart.Format("2006-01-02")
		}
		if endTime == "" {
			endTime = defaultEnd.Format("2006-01-02")
		}
	}

	startDate, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start date: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end date: %v", err)
	}

	timePeriod.From = to.Ptr(startDate)
	timePeriod.To = to.Ptr(endDate)

	// Daily granularity
	azureGranularity := getGranularityFromString("DAILY")

	// Resource group name grouping
	groupBy := &armcostmanagement.QueryGrouping{
		Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
		Name: to.Ptr("ResourceGroupName"),
	}

	// Build filter for resource_group if provided
	filter := buildFilterExpression(d, "ResourceGroupName")

	return &AzureCostQueryInput{
		Timeframe:   timeframe,
		Granularity: azureGranularity,
		GroupBy:     groupBy,
		Scope:       "/subscriptions/" + subscriptionID,
		TimePeriod:  timePeriod,
		Filter:      filter,
	}, nil
}
