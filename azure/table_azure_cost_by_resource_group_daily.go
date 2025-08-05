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
	queryDef, scope, err := buildCostByResourceGroupDailyInput(ctx, d)
	if err != nil {
		return nil, err
	}
	return streamCostAndUsage(ctx, d, queryDef, scope, "ResourceGroupName")
}

func buildCostByResourceGroupDailyInput(ctx context.Context, d *plugin.QueryData) (armcostmanagement.QueryDefinition, string, error) {
	// Get scope from quals, default to placeholder if not provided
	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "/subscriptions/placeholder" // Will be resolved in streamCostAndUsage
	}

	// Get cost type from quals, default to ActualCost
	costType := d.EqualsQualString("type")
	if costType == "" {
		costType = "ActualCost"
	}

	// Set timeframe and time period using new usage_date logic
	timeframe := armcostmanagement.TimeframeTypeCustom
	timePeriod := &armcostmanagement.QueryTimePeriod{}

	// Get time range from period_start/period_end quals
	startTime, endTime := getPeriodTimeRange(d, "DAILY")

	startDate, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return armcostmanagement.QueryDefinition{}, "", fmt.Errorf("failed to parse start date: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		return armcostmanagement.QueryDefinition{}, "", fmt.Errorf("failed to parse end date: %v", err)
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

	// Get dynamic columns based on query context
	_ = getColumnsFromQueryContext(d.QueryContext) // Not used in grouped queries

	// Build aggregation based on requested columns
	aggregation := make(map[string]*armcostmanagement.QueryAggregation)

	// Determine which metrics to include
	metrics := getMetricsByQueryContext(d.QueryContext)
	if len(metrics) == 0 {
		// Default metrics if none specified (only cost metrics)
		metrics = []string{"PreTaxCost"}
	}

	// Add aggregations (Azure limit is 2)
	for i, metric := range metrics {
		if i >= 2 {
			break
		}
		aggregation[metric] = &armcostmanagement.QueryAggregation{
			Function: to.Ptr(armcostmanagement.FunctionTypeSum),
			Name:     to.Ptr(metric),
		}
	}

	// Create dataset with grouping
	dataset := &armcostmanagement.QueryDataset{
		Granularity: &azureGranularity,
		Aggregation: aggregation,
		Grouping:    []*armcostmanagement.QueryGrouping{groupBy},
	}

	// Add filter if specified
	if filter != nil {
		dataset.Filter = filter
	}

	// Create QueryDefinition
	queryDef := armcostmanagement.QueryDefinition{
		Type:      to.Ptr(getCostTypeFromString(costType)),
		Timeframe: to.Ptr(timeframe),
		Dataset:   dataset,
	}

	// Set TimePeriod if using Custom timeframe
	if timeframe == armcostmanagement.TimeframeTypeCustom && timePeriod != nil {
		queryDef.TimePeriod = timePeriod
	}

	return queryDef, scope, nil
}
