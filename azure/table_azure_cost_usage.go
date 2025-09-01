package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureCostUsage(_ context.Context) *plugin.Table {
	keyColumns := costManagementKeyColumns()
	keyColumns = append(
		plugin.KeyColumnSlice{
			{
				Name:    "granularity",
				Require: plugin.Required,
			},
			{
				Name:    "dimension_type_1",
				Require: plugin.AnyOf,
			},
			{
				Name:    "dimension_type_2",
				Require: plugin.AnyOf,
			},
			{
				Name:    "dimension_types",
				Require: plugin.AnyOf,
			},
		}, keyColumns...,
	)

	return &plugin.Table{
		Name:        "azure_cost_usage",
		Description: "Azure Cost Management - Cost and Usage with flexible dimensions",
		List: &plugin.ListConfig{
			KeyColumns: keyColumns,
			Hydrate:    listCostUsage,
			Tags:       map[string]string{"service": "Microsoft.CostManagement", "action": "Query"},
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
				{
					Name:        "dimensions",
					Description: "The dimensions values in the form of key value pairs.",
					Type:        proto.ColumnType_JSON,
					Transform:   transform.FromField("Dimensions"),
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
				{
					Name:        "dimension_types",
					Description: "The dimension keys in the form of array of strings.",
					Type:        proto.ColumnType_JSON,
					Hydrate:     hydrateCostUsageQuals,
				},
			}),
		),
	}
}

//// LIST FUNCTION

func listCostUsage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	queryDef, scope, dim1, dim2, dims, err := buildCostUsageInputFromQuals(ctx, d)
	if err != nil {
		return nil, err
	}
	if len(dims) > 0 {
		return streamCostAndUsage(ctx, d, queryDef, scope, dims...)
	}
	return streamCostAndUsage(ctx, d, queryDef, scope, dim1, dim2)
}

func buildCostUsageInputFromQuals(ctx context.Context, d *plugin.QueryData) (armcostmanagement.QueryDefinition, string, string, string, []string, error) {
	granularity := strings.ToUpper(d.EqualsQuals["granularity"].GetStringValue())

	// Get scope from quals, default to placeholder if not provided
	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "/subscriptions/placeholder" // Will be resolved in streamCostAndUsage
	}

	// Get cost type from quals, default to ActualCost
	costType := d.EqualsQualString("type")
	if costType == "" {
		return armcostmanagement.QueryDefinition{}, "", "", "", []string{}, fmt.Errorf("missing required qual 'type' (ActualCost | AmortizedCost)")
	}

	// Set timeframe and time period
	timeframe := armcostmanagement.TimeframeTypeCustom
	timePeriod := &armcostmanagement.QueryTimePeriod{}

	// Get time range from period_start/period_end quals
	startTime, endTime := getPeriodTimeRange(d, granularity)

	// Set default time range if no quals provided
	if startTime == "" || endTime == "" {
		var defaultStart, defaultEnd time.Time
		switch granularity {
		case "MONTHLY", "DAILY":
			// Default: 1 year back
			defaultEnd = time.Now()
			defaultStart = defaultEnd.AddDate(0, -11, -30)
		default:
			// Default: Last 30 days
			defaultEnd = time.Now()
			defaultStart = defaultEnd.AddDate(0, 0, -30)
		}

		if startTime == "" {
			startTime = defaultStart.Format("2006-01-02")
		}
		if endTime == "" {
			endTime = defaultEnd.Format("2006-01-02")
		}
	}

	startDate, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return armcostmanagement.QueryDefinition{}, "", "", "", []string{}, fmt.Errorf("failed to parse start date: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		return armcostmanagement.QueryDefinition{}, "", "", "", []string{}, fmt.Errorf("failed to parse end date: %v", err)
	}

	timePeriod.From = to.Ptr(startDate)
	timePeriod.To = to.Ptr(endDate)

	azureGranularity := getGranularityFromString(granularity)

	// Get dimensions
	dim1 := d.EqualsQuals["dimension_type_1"].GetStringValue()
	dim2 := d.EqualsQuals["dimension_type_2"].GetStringValue()
	dims := d.EqualsQuals["dimension_types"].GetJsonbValue()

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

	var dimensions []string
	if dims != "" {
		err := json.Unmarshal([]byte(dims), &dimensions)
		if err != nil {
			return armcostmanagement.QueryDefinition{}, "", "", "", []string{}, fmt.Errorf("failed to parse dimensions: %v", err)
		}
		for _, dim := range dimensions {
			groupings = append(groupings, &armcostmanagement.QueryGrouping{
				Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
				Name: to.Ptr(dim),
			})
		}
	}

	// Get dynamic columns based on query context
	requestColumns := getColumnsFromQueryContext(d.QueryContext)

	// Build aggregation based on requested columns
	aggregation := make(map[string]*armcostmanagement.QueryAggregation)

	// Determine which metrics to include (from global CostMetrics)
	for i, metric := range CostMetrics {
		if i >= 2 { // Azure allows max 2 aggregations
			break
		}
		aggregation[metric] = &armcostmanagement.QueryAggregation{
			Function: to.Ptr(armcostmanagement.FunctionTypeSum),
			Name:     to.Ptr(metric),
		}
	}

	// Azure API restriction: Cannot use Configuration with Grouping/Aggregation
	// If we have grouping, use traditional approach without Configuration
	// If no grouping, use Configuration approach for column selection

	var dataset *armcostmanagement.QueryDataset

	if len(groupings) > 0 {
		// Traditional approach with grouping - no Configuration
		dataset = &armcostmanagement.QueryDataset{
			Granularity: &azureGranularity,
			Aggregation: aggregation,
		}

		// Add grouping if specified - Azure supports up to 2 groupings
		if len(groupings) > 0 {
			dataset.Grouping = groupings
		}
	} else {
		// Configuration approach without grouping - for flexible column selection
		dataset = &armcostmanagement.QueryDataset{
			Granularity: &azureGranularity,
			Configuration: &armcostmanagement.QueryDatasetConfiguration{
				Columns: requestColumns,
			},
		}
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

	return queryDef, scope, dim1, dim2, dimensions, nil
}

//// HYDRATE FUNCTIONS

func hydrateCostUsageQuals(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dimensionTypesJson := d.EqualsQuals["dimension_types"].GetJsonbValue()
	var dimensionTypes []string
	if dimensionTypesJson != "" {
		err := json.Unmarshal([]byte(dimensionTypesJson), &dimensionTypes)
		if err != nil {
			return nil, err
		}
	}
	return &AzureCostUsageQuals{
		Granularity:    d.EqualsQuals["granularity"].GetStringValue(),
		DimensionType1: d.EqualsQuals["dimension_type_1"].GetStringValue(),
		DimensionType2: d.EqualsQuals["dimension_type_2"].GetStringValue(),
		DimensionTypes: dimensionTypes,
	}, nil
}

type AzureCostUsageQuals struct {
	Granularity    string
	DimensionType1 string
	DimensionType2 string
	DimensionTypes []string
}
