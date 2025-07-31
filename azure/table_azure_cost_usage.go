package azure

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAzureCostUsage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cost_usage",
		Description: "Azure Cost Management - Cost and Usage with flexible dimensions",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "granularity",
					Require: plugin.Required,
				},
				{
					Name:    "dimension_type_1",
					Require: plugin.Required,
				},
				{
					Name:    "dimension_type_2",
					Require: plugin.Required,
				},
				{
					Name:       "period_start",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: query_cache.CacheMatchExact,
				},
				{
					Name:       "period_end",
					Require:    plugin.Optional,
					Operators:  []string{">", ">=", "=", "<", "<="},
					CacheMatch: query_cache.CacheMatchExact,
				},
			},
			Hydrate: listCostUsage,
			Tags:    map[string]string{"service": "Microsoft.CostManagement", "action": "Query"},
		},
		Columns: azureColumns(
			costManagementColumns([]*plugin.Column{
				{
					Name:        "dimension_1",
					Description: "The first dimension value. Valid dimension types include ResourceGroupName, ServiceName, Location, ResourceType, MeterCategory, MeterSubCategory, MeterName, ResourceId, PublisherType, ChargeType, ReservationId, Frequency",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension1"),
				},
				{
					Name:        "dimension_2",
					Description: "The second dimension value. Valid dimension types include ResourceGroupName, ServiceName, Location, ResourceType, MeterCategory, MeterSubCategory, MeterName, ResourceId, PublisherType, ChargeType, ReservationId, Frequency",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("Dimension2"),
				},
				// Quals columns - to filter the lookups
				{
					Name:        "granularity",
					Description: "The Azure cost granularity. Valid values are Daily or Monthly.",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostUsageQuals,
				},
				{
					Name:        "dimension_type_1",
					Description: "The first dimension to group results by. Valid values include ResourceGroupName, ServiceName, Location, ResourceType, MeterCategory, MeterSubCategory, MeterName, ResourceId, PublisherType, ChargeType, ReservationId, Frequency",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostUsageQuals,
				},
				{
					Name:        "dimension_type_2",
					Description: "The second dimension to group results by. Valid values include ResourceGroupName, ServiceName, Location, ResourceType, MeterCategory, MeterSubCategory, MeterName, ResourceId, PublisherType, ChargeType, ReservationId, Frequency",
					Type:        proto.ColumnType_STRING,
					Hydrate:     hydrateCostUsageQuals,
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostUsage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildCostUsageInputFromQuals(ctx, d)
	return streamCostAndUsage(ctx, d, params)
}

func buildCostUsageInputFromQuals(ctx context.Context, keyQuals *plugin.QueryData) *AzureCostQueryInput {
	granularity := strings.ToUpper(keyQuals.EqualsQuals["granularity"].GetStringValue())

	// Get subscription ID
	subscriptionID := keyQuals.EqualsQualString("subscription_id")
	if subscriptionID == "" {
		subscriptionID = "placeholder" // Will be replaced in streamCostAndUsage
	}

	// Set timeframe and time period
	var timeframe armcostmanagement.TimeframeType
	var timePeriod *armcostmanagement.QueryTimePeriod = nil

	// Get time range from quals
	startTime, endTime := getCostUsageTimeRange(keyQuals, granularity)

	// Check if user provided specific time range
	hasTimeFilter := false
	if quals := keyQuals.Quals["period_start"]; quals != nil && len(quals.Quals) > 0 {
		hasTimeFilter = true
	}
	if quals := keyQuals.Quals["period_end"]; quals != nil && len(quals.Quals) > 0 {
		hasTimeFilter = true
	}

	if hasTimeFilter {
		timeframe = armcostmanagement.TimeframeTypeCustom
		// Parse time strings to time.Time
		startDate, _ := time.Parse("2006-01-02", startTime)
		endDate, _ := time.Parse("2006-01-02", endTime)
		timePeriod = &armcostmanagement.QueryTimePeriod{
			From: to.Ptr(startDate),
			To:   to.Ptr(endDate),
		}
	} else {
		// Use MonthToDate for monthly aggregation, or default timeframe
		if granularity == "MONTHLY" {
			timeframe = armcostmanagement.TimeframeTypeMonthToDate
		} else {
			timeframe = armcostmanagement.TimeframeTypeMonthToDate // Default
		}
	}

	azureGranularity := getGranularityFromString(granularity)

	// Get dimensions
	dim1 := keyQuals.EqualsQuals["dimension_type_1"].GetStringValue()
	dim2 := keyQuals.EqualsQuals["dimension_type_2"].GetStringValue()

	// Build GroupBy - Azure supports up to 2 grouping dimensions
	var groupings []*armcostmanagement.QueryGrouping
	if dim1 != "" {
		groupings = append(groupings, &armcostmanagement.QueryGrouping{
			Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
			Name: to.Ptr(dim1),
		})
	}
	if dim2 != "" {
		groupings = append(groupings, &armcostmanagement.QueryGrouping{
			Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
			Name: to.Ptr(dim2),
		})
	}

	// Create input parameters
	params := &AzureCostQueryInput{
		Timeframe:   timeframe,
		Granularity: azureGranularity,
		Scope:       "/subscriptions/" + subscriptionID,
		TimePeriod:  timePeriod,
		Filter:      nil, // No dimension-specific filters for this table
	}

	// Set groupings - use both GroupBy and GroupBy2 for multiple dimensions
	if len(groupings) > 0 {
		params.GroupBy = groupings[0]
	}
	if len(groupings) > 1 {
		params.GroupBy2 = groupings[1]
	}

	return params
}

func getCostUsageTimeRange(keyQuals *plugin.QueryData, granularity string) (string, string) {
	timeFormat := "2006-01-02"

	// Default time range based on granularity
	var defaultStart time.Time
	var defaultEnd time.Time

	switch granularity {
	case "MONTHLY":
		// Default: 1 year back
		defaultEnd = time.Now()
		defaultStart = defaultEnd.AddDate(0, -11, -30)
	case "DAILY":
		// Default: Last 30 days
		defaultEnd = time.Now()
		defaultStart = defaultEnd.AddDate(0, 0, -30)
	default:
		// Default: Last 30 days
		defaultEnd = time.Now()
		defaultStart = defaultEnd.AddDate(0, 0, -30)
	}

	startTime := defaultStart.Format(timeFormat)
	endTime := defaultEnd.Format(timeFormat)

	// Process period_start quals (similar to AWS)
	if keyQuals.Quals["period_start"] != nil && len(keyQuals.Quals["period_start"].Quals) <= 1 {
		for _, q := range keyQuals.Quals["period_start"].Quals {
			t := q.Value.GetTimestampValue().AsTime().Format(timeFormat)
			switch q.Operator {
			case "=", ">=", ">":
				startTime = t
			case "<", "<=":
				endTime = t
			}
		}
	}

	// Process period_end quals (similar to AWS)
	if keyQuals.Quals["period_end"] != nil && len(keyQuals.Quals["period_end"].Quals) <= 1 {
		for _, q := range keyQuals.Quals["period_end"].Quals {
			t := q.Value.GetTimestampValue().AsTime().Format(timeFormat)
			switch q.Operator {
			case "=", ">=", ">":
				if startTime == defaultStart.Format(timeFormat) {
					startTime = t
				}
			case "<", "<=":
				if endTime == defaultEnd.Format(timeFormat) {
					endTime = t
				}
			}
		}
	}

	return startTime, endTime
}

//// HYDRATE FUNCTIONS

func hydrateCostUsageQuals(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return &AzureCostUsageQuals{
		Granularity:    d.EqualsQuals["granularity"].GetStringValue(),
		DimensionType1: d.EqualsQuals["dimension_type_1"].GetStringValue(),
		DimensionType2: d.EqualsQuals["dimension_type_2"].GetStringValue(),
	}, nil
}

type AzureCostUsageQuals struct {
	Granularity    string
	DimensionType1 string
	DimensionType2 string
}
