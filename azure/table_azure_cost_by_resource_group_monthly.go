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
	granularity := "MONTHLY"
	params, err := buildCostByResourceGroupInput(ctx, granularity, d)
	if err != nil {
		return nil, err
	}
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByResourceGroupInput(ctx context.Context, granularity string, d *plugin.QueryData) (*AzureCostQueryInput, error) {
	// Get subscription ID
	subscriptionID := d.EqualsQualString("subscription_id")
	if subscriptionID == "" {
		// We'll get subscription ID during execution using hydrate data
		subscriptionID = "placeholder" // Will be replaced in streamCostAndUsage
	}

	// Set timeframe and granularity to match working raw API call
	var timePeriod *armcostmanagement.QueryTimePeriod

	// Get time range from usage_date quals using simplified approach
	startTime, endTime := getUsageDateTimeRange(d, granularity)

	// Set default time range if no quals provided
	if startTime == "" || endTime == "" {
		defaultEnd := time.Now()
		defaultStart := defaultEnd.AddDate(0, -11, -30) // 1 year back for monthly
		if startTime == "" {
			startTime = defaultStart.Format("2006-01-02")
		}
		if endTime == "" {
			endTime = defaultEnd.Format("2006-01-02")
		}
	}

	// Parse time strings to time.Time
	startDate, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start date: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end date: %v", err)
	}
	timePeriod = &armcostmanagement.QueryTimePeriod{
		From: to.Ptr(startDate),
		To:   to.Ptr(endDate),
	}

	azureGranularity := getGranularityFromString(granularity) // Use Monthly granularity when available

	// Build GroupBy for ResourceGroup
	groupBy := &armcostmanagement.QueryGrouping{
		Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
		Name: to.Ptr("ResourceGroupName"), // Use correct Azure API dimension name
	}

	// Build filter expressions from quals
	filter := buildFilterExpression(d, "ResourceGroupName")

	// Create input parameters
	params := &AzureCostQueryInput{
		Timeframe:   armcostmanagement.TimeframeTypeCustom,
		Granularity: azureGranularity,
		GroupBy:     groupBy,
		Scope:       "/subscriptions/" + subscriptionID,
		TimePeriod:  timePeriod,
		Filter:      filter,
	}

	return params, nil
}
