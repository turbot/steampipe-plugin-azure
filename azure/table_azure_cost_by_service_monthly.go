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
	granularity := "MONTHLY"
	params, err := buildCostByServiceInput(ctx, granularity, d)
	if err != nil {
		return nil, err
	}
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByServiceInput(ctx context.Context, granularity string, d *plugin.QueryData) (*AzureCostQueryInput, error) {
	// Get subscription ID
	subscriptionID := d.EqualsQualString("subscription_id")
	if subscriptionID == "" {
		// We'll get subscription ID during execution using hydrate data
		subscriptionID = "placeholder" // Will be replaced in streamCostAndUsage
	}

	// Set timeframe and granularity to match working raw API call
	timePeriod := &armcostmanagement.QueryTimePeriod{}

	// Get time range from quals
	startTime, endTime := getCostUsageTimeRange(d, granularity)

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

	// Build GroupBy for ServiceName
	groupBy := &armcostmanagement.QueryGrouping{
		Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
		Name: to.Ptr("ServiceName"),
	}

	// Build filter expressions from quals
	filter := buildFilterExpression(d, "ServiceName")

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
