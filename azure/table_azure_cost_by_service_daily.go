package azure

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureCostByServiceDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cost_by_service_daily",
		Description: "Azure Cost Management - Daily cost by service",
		List: &plugin.ListConfig{
			KeyColumns: append(costManagementKeyColumns(),
				&plugin.KeyColumn{
					Name:    "service_name",
					Require: plugin.Optional,
				},
			),
			Hydrate: listCostByServiceDaily,
			Tags:    map[string]string{"service": "Microsoft.CostManagement", "action": "Query"},
		},
		Columns: azureColumns(
			costManagementColumns([]*plugin.Column{
				{
					Name:        "service_name",
					Description: "The name of the Azure service (e.g., Virtual Machines, Storage, etc.)",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostByServiceDaily(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostByServiceDailyInput(d)
	return streamCostAndUsage(ctx, d, params)
}

func buildCostByServiceDailyInput(d *plugin.QueryData) *AzureCostQueryInput {
	// Get time range from quals with daily defaults
	startTime, endTime := getTimeRangeFromQuals(d, "DAILY")

	// Get subscription ID (will be handled in streamCostAndUsage if empty)
	subscriptionID := d.EqualsQualString("subscription_id")
	if subscriptionID == "" {
		subscriptionID = "placeholder"
	}

	// Set timeframe and time period
	var timeframe armcostmanagement.TimeframeType
	var timePeriod *armcostmanagement.QueryTimePeriod = nil

	// Check if user provided specific time range
	hasTimeFilter := false
	if quals := d.Quals["period_start"]; quals != nil && len(quals.Quals) > 0 {
		hasTimeFilter = true
	}
	if quals := d.Quals["period_end"]; quals != nil && len(quals.Quals) > 0 {
		hasTimeFilter = true
	}

	if hasTimeFilter {
		timeframe = armcostmanagement.TimeframeTypeCustom
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
	} else {
		// Always use Custom timeframe for daily data with last 7 days
		timeframe = armcostmanagement.TimeframeTypeCustom
		endDate, err := time.Parse("2006-01-02", endTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		startDate, err := time.Parse("2006-01-02", startTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %v", err)
		}
		timePeriod = &armcostmanagement.QueryTimePeriod{
			From: to.Ptr(startDate),
			To:   to.Ptr(endDate),
		}
	}

	// Daily granularity
	azureGranularity := getGranularityFromString("DAILY")

	// Service name grouping
	groupBy := &armcostmanagement.QueryGrouping{
		Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
		Name: to.Ptr("ServiceName"),
	}

	// Build filter for service_name if provided
	filter := buildFilterExpression(d, "ServiceName")

	return &AzureCostQueryInput{
		Timeframe:   timeframe,
		Granularity: azureGranularity,
		GroupBy:     groupBy,
		Scope:       "/subscriptions/" + subscriptionID,
		TimePeriod:  timePeriod,
		Filter:      filter,
	}
}
