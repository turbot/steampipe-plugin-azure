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

// CostManagementRow represents a flattened cost management result row with all cost types (like AWS)
type CostManagementRow struct {
	PeriodStart *string
	PeriodEnd   *string

	// Dimension values (populated based on GroupBy)
	Dimension1 *string // Generic dimension field (could be ResourceGroup, ServiceName, etc.)
	Dimension2 *string // Second dimension field for multi-dimensional grouping

	// Cost metrics (following AWS naming conventions)
	// Actual costs
	UnblendedCostAmount *float64 // Primary actual cost
	UnblendedCostUnit   *string
	PreTaxCostAmount    *float64 // Pre-tax cost
	PreTaxCostUnit      *string

	// Amortized costs (for reservations)
	AmortizedCostAmount *float64
	AmortizedCostUnit   *string

	// Usage metrics
	UsageQuantityAmount *float64
	UsageQuantityUnit   *string

	// Metadata
	Estimated *bool
	Currency  *string

	// Common properties
	SubscriptionID   *string
	SubscriptionName *string
}

// AzureCostQueryInput contains all parameters needed for Azure cost queries (similar to AWS GetCostAndUsageInput)
type AzureCostQueryInput struct {
	Timeframe   armcostmanagement.TimeframeType
	Granularity armcostmanagement.GranularityType
	GroupBy     *armcostmanagement.QueryGrouping   // Optional grouping (primary)
	GroupBy2    *armcostmanagement.QueryGrouping   // Optional second grouping
	Metrics     []string                           // List of metrics to fetch
	TimePeriod  *armcostmanagement.QueryTimePeriod // Optional custom time period
	Scope       string                             // Subscription scope
	Filter      *armcostmanagement.QueryFilter     // Optional filter expressions
}

// getGranularityFromString converts granularity string to Azure GranularityType
// Azure supports both Daily and Monthly, but we need to use the right type
func getGranularityFromString(granularity string) armcostmanagement.GranularityType {
	switch granularity {
	case "MONTHLY":
		// Use Monthly granularity directly as string (same as raw API)
		return armcostmanagement.GranularityType("Monthly")
	case "DAILY":
		return armcostmanagement.GranularityTypeDaily
	default:
		return armcostmanagement.GranularityTypeDaily
	}
}

// AllCostMetrics returns all available cost metrics for Azure Cost Management (like AWS)
func AllCostMetrics() []string {
	return []string{
		"PreTaxCost",
		"Cost",
		"UsageQuantity",
	}
}

// getMetricsByQueryContext dynamically determines which metrics to fetch based on query columns (like AWS)
func getMetricsByQueryContext(qc *plugin.QueryContext) []string {
	queryColumns := qc.Columns
	var metrics []string

	for _, c := range queryColumns {
		switch c {
		case "unblended_cost_amount", "unblended_cost_unit":
			metrics = append(metrics, "PreTaxCost")
		case "pre_tax_cost_amount", "pre_tax_cost_unit":
			metrics = append(metrics, "PreTaxCost")
		case "amortized_cost_amount", "amortized_cost_unit":
			metrics = append(metrics, "PreTaxCost") // Use PreTaxCost for amortized as well
		case "usage_quantity_amount", "usage_quantity_unit":
			metrics = append(metrics, "UsageQuantity")
		}
	}

	return removeDuplicates(metrics)
}

// getColumnsFromQueryContext determines which columns to request from Azure API
func getColumnsFromQueryContext(qc *plugin.QueryContext) []*string {
	queryColumns := qc.Columns
	var columns []*string

	// Always include basic columns
	columns = append(columns, to.Ptr("Currency"))

	// Add columns based on what user is querying
	needsCost := false
	needsUsage := false

	for _, c := range queryColumns {
		switch c {
		case "unblended_cost_amount", "unblended_cost_unit", "pre_tax_cost_amount", "pre_tax_cost_unit", "amortized_cost_amount", "amortized_cost_unit":
			needsCost = true
		case "usage_quantity_amount", "usage_quantity_unit":
			needsUsage = true
		}
	}

	if needsCost {
		columns = append(columns, to.Ptr("PreTaxCost"))
	}
	if needsUsage {
		columns = append(columns, to.Ptr("UsageQuantity"))
	}

	// If no specific columns requested, include both
	if !needsCost && !needsUsage {
		columns = append(columns, to.Ptr("PreTaxCost"), to.Ptr("UsageQuantity"))
	}

	return columns
}

// removeDuplicates removes duplicate metrics from slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// getTimeRangeFromQuals extracts start and end dates from period_start and period_end quals
func getTimeRangeFromQuals(d *plugin.QueryData, granularity string) (string, string) {
	timeFormat := "2006-01-02"

	// Default time range based on granularity
	var defaultStart time.Time
	var defaultEnd time.Time

	switch granularity {
	case "MONTHLY":
		// Default: 1 year back (Azure only supports Daily granularity, but we aggregate over longer periods for monthly data)
		defaultEnd = time.Now()
		defaultStart = defaultEnd.AddDate(0, -11, -30) // 1 year back
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

	// Process period_start quals
	if quals := d.Quals["period_start"]; quals != nil {
		for _, qual := range quals.Quals {
			if qual.Value != nil && qual.Value.GetTimestampValue() != nil {
				qualTime := qual.Value.GetTimestampValue().AsTime()
				switch qual.Operator {
				case "=":
					startTime = qualTime.Format(timeFormat)
				case ">=", ">":
					if qual.Operator == ">" {
						qualTime = qualTime.AddDate(0, 0, 1) // Next day for > operator
					}
					startTime = qualTime.Format(timeFormat)
				case "<=", "<":
					if qual.Operator == "<" {
						qualTime = qualTime.AddDate(0, 0, -1) // Previous day for < operator
					}
					endTime = qualTime.Format(timeFormat)
				}
			}
		}
	}

	// Process period_end quals
	if quals := d.Quals["period_end"]; quals != nil {
		for _, qual := range quals.Quals {
			if qual.Value != nil && qual.Value.GetTimestampValue() != nil {
				qualTime := qual.Value.GetTimestampValue().AsTime()
				switch qual.Operator {
				case "=":
					endTime = qualTime.Format(timeFormat)
				case "<=", "<":
					if qual.Operator == "<" {
						qualTime = qualTime.AddDate(0, 0, -1) // Previous day for < operator
					}
					endTime = qualTime.Format(timeFormat)
				case ">=", ">":
					if qual.Operator == ">" {
						qualTime = qualTime.AddDate(0, 0, 1) // Next day for > operator
					}
					startTime = qualTime.Format(timeFormat)
				}
			}
		}
	}

	return startTime, endTime
}

// buildFilterExpression creates filter expressions from table key column quals
func buildFilterExpression(d *plugin.QueryData, dimensionName string) *armcostmanagement.QueryFilter {
	var filters []*armcostmanagement.QueryFilter
	// Process dimension-specific quals (like service_name = 'Storage')
	for _, keyQual := range d.Table.List.KeyColumns {
		if keyQual.Name == "period_start" || keyQual.Name == "period_end" || keyQual.Name == "subscription_id" {
			continue // Skip time and subscription quals
		}

		filterQual := d.Quals[keyQual.Name]
		if filterQual == nil {
			continue
		}

		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				value := qual.Value.GetStringValue()
				if value == "" {
					continue
				}

				switch qual.Operator {
				case "=":
					filter := &armcostmanagement.QueryFilter{
						Dimensions: &armcostmanagement.QueryComparisonExpression{
							Name:     to.Ptr(dimensionName),
							Operator: to.Ptr(armcostmanagement.QueryOperatorTypeIn),
							Values:   []*string{to.Ptr(value)},
						},
					}
					filters = append(filters, filter)
				case "<>":
					// Azure doesn't have a Not field, so we'll use a separate filter for exclusions
					// For now, we'll skip <> operators or handle them differently
					// TODO: Implement exclusion logic using different approach
					continue
				}
			}
		}
	}

	// Combine filters with AND logic
	if len(filters) == 0 {
		return nil
	} else if len(filters) == 1 {
		return filters[0]
	} else {
		return &armcostmanagement.QueryFilter{
			And: filters,
		}
	}
}

// costManagementColumns returns the standard cost management columns (like AWS costExplorerColumns)
func costManagementColumns(columns []*plugin.Column) []*plugin.Column {
	standardColumns := []*plugin.Column{
		{
			Name:        "period_start",
			Description: "Start timestamp for this cost metric.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("PeriodStart"),
		},
		{
			Name:        "period_end",
			Description: "End timestamp for this cost metric.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("PeriodEnd"),
		},
		{
			Name:        "estimated",
			Description: "Whether the cost data is estimated.",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Estimated"),
		},
		// Actual costs (similar to AWS unblended costs)
		{
			Name:        "unblended_cost_amount",
			Description: "Unblended costs represent your usage costs on the day they are charged to you. In finance terms, they represent your costs on a cash basis of accounting.",
			Type:        proto.ColumnType_DOUBLE,
			Transform:   transform.FromField("UnblendedCostAmount"),
		},
		{
			Name:        "unblended_cost_unit",
			Description: "Unit type for unblended costs.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UnblendedCostUnit"),
		},
		{
			Name:        "pre_tax_cost_amount",
			Description: "Pre-tax cost amount for the period.",
			Type:        proto.ColumnType_DOUBLE,
			Transform:   transform.FromField("PreTaxCostAmount"),
		},
		{
			Name:        "pre_tax_cost_unit",
			Description: "Unit type for pre-tax costs.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("PreTaxCostUnit"),
		},
		// Amortized costs (for reservations)
		{
			Name:        "amortized_cost_amount",
			Description: "This cost metric reflects the effective cost of the upfront and monthly reservation fees spread across the billing period. By default, Cost Explorer shows the fees for Reserved Instances as a spike on the day that you're charged, but if you choose to show costs as amortized costs, the costs are amortized over the billing period. This means that the costs are broken out into the effective daily rate. Azure estimates your amortized costs by combining your unblended costs with the amortized portion of your upfront and recurring reservation fees.",
			Type:        proto.ColumnType_DOUBLE,
			Transform:   transform.FromField("AmortizedCostAmount"),
		},
		{
			Name:        "amortized_cost_unit",
			Description: "Unit type for amortized costs.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("AmortizedCostUnit"),
		},
		// Usage metrics
		{
			Name:        "usage_quantity_amount",
			Description: "The amount of usage that you incurred. NOTE: If you return the UsageQuantity metric, the service aggregates all usage numbers without taking into account the units. For example, if you aggregate usageQuantity across all of Azure Compute, the results aren't meaningful because Azure Compute hours and data transfer are measured in different units (for example, hours vs. GB).",
			Type:        proto.ColumnType_DOUBLE,
			Transform:   transform.FromField("UsageQuantityAmount"),
		},
		{
			Name:        "usage_quantity_unit",
			Description: "Unit type for usage quantity.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UsageQuantityUnit"),
		},
	}

	// Prepend table-specific columns to standard columns
	return append(columns, standardColumns...)
}

// costManagementKeyColumns returns the standard key columns for cost management tables (like AWS)
func costManagementKeyColumns() plugin.KeyColumnSlice {
	return plugin.KeyColumnSlice{
		{
			Name:      "subscription_id",
			Require:   plugin.Optional,
			Operators: []string{"="},
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
	}
}

// getCostManagementClient creates a new cost management client
func getCostManagementClient(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*armcostmanagement.QueryClient, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	client, err := armcostmanagement.NewQueryClient(session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// streamCostAndUsage is the generic function for streaming cost data (like AWS)
func streamCostAndUsage(ctx context.Context, d *plugin.QueryData, params *AzureCostQueryInput) (interface{}, error) {
	client, err := getCostManagementClient(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	// Resolve subscription ID if it's a placeholder
	if params.Scope == "/subscriptions/placeholder" {
		subscriptionData, err := getSubscriptionID(ctx, d, nil)
		if err != nil {
			return nil, err
		}
		subscriptionID := subscriptionData.(string)
		params.Scope = "/subscriptions/" + subscriptionID
	}

	// Get dynamic columns based on query context
	requestedColumns := getColumnsFromQueryContext(d.QueryContext)

	// Build aggregation based on requested columns
	aggregation := make(map[string]*armcostmanagement.QueryAggregation)

	// Determine which metrics to include
	metrics := getMetricsByQueryContext(d.QueryContext)
	if len(metrics) == 0 {
		// Default metrics if none specified
		metrics = []string{"PreTaxCost", "UsageQuantity"}
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

	// Azure API restriction: Cannot use Configuration with Grouping/Aggregation
	// If we have grouping, use traditional approach without Configuration
	// If no grouping, use Configuration approach for column selection

	var dataset *armcostmanagement.QueryDataset

	if params.GroupBy != nil || params.GroupBy2 != nil {
		// Traditional approach with grouping - no Configuration
		dataset = &armcostmanagement.QueryDataset{
			Granularity: &params.Granularity,
			Aggregation: aggregation,
		}

		// Add grouping if specified - Azure supports up to 2 groupings
		var groupings []*armcostmanagement.QueryGrouping
		if params.GroupBy != nil {
			groupings = append(groupings, params.GroupBy)
		}
		if params.GroupBy2 != nil {
			groupings = append(groupings, params.GroupBy2)
		}
		if len(groupings) > 0 {
			dataset.Grouping = groupings
		}
	} else {
		// Configuration approach without grouping - for flexible column selection
		dataset = &armcostmanagement.QueryDataset{
			Granularity: &params.Granularity,
			Configuration: &armcostmanagement.QueryDatasetConfiguration{
				Columns: requestedColumns,
			},
		}
	}

	// Add filter if specified
	if params.Filter != nil {
		dataset.Filter = params.Filter
	}

	// Build query definition - using Usage type like the working example
	queryDef := armcostmanagement.QueryDefinition{
		Type:      to.Ptr(armcostmanagement.ExportTypeUsage),
		Timeframe: to.Ptr(params.Timeframe),
		Dataset:   dataset,
	}

	// Set TimePeriod only if using Custom timeframe
	if params.Timeframe == armcostmanagement.TimeframeTypeCustom && params.TimePeriod != nil {
		queryDef.TimePeriod = params.TimePeriod
	}

	// Execute query
	result, err := client.Usage(ctx, params.Scope, queryDef, nil)
	if err != nil {
		return nil, err
	}

	// Process results - simplified
	rowMap := make(map[string]*CostManagementRow)
	processQueryResults(&result.QueryResult, params, "actual", rowMap)

	// Stream results
	for _, row := range rowMap {
		d.StreamListItem(ctx, *row)

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// processQueryResults processes query results and merges them into the row map
func processQueryResults(result *armcostmanagement.QueryResult, params *AzureCostQueryInput, costType string, rowMap map[string]*CostManagementRow) {
	if result.Properties == nil || result.Properties.Columns == nil || result.Properties.Rows == nil {
		return
	}

	// Extract column mapping for easier access
	columnMap := make(map[string]int)
	for i, col := range result.Properties.Columns {
		if col.Name != nil {
			columnMap[*col.Name] = i
		}
	}

	// Process each row
	for _, row := range result.Properties.Rows {
		if len(row) == 0 {
			continue
		}

		// Build row key - date should always be first
		keyParts := []string{}

		// Identify the date column based on granularity or column names
		var dateColumnName string
		var dateIdx int = -1

		// For daily granularity, look for UsageDate column
		if params.Granularity == armcostmanagement.GranularityTypeDaily {
			if idx, ok := columnMap["UsageDate"]; ok {
				dateColumnName = "UsageDate"
				dateIdx = idx
			}
		} else {
			// For monthly granularity, look for BillingMonth column
			if idx, ok := columnMap["BillingMonth"]; ok {
				dateColumnName = "BillingMonth"
				dateIdx = idx
			}
		}

		// Get date value
		var dateStr string
		if dateIdx != -1 && len(row) > dateIdx && row[dateIdx] != nil {
			if date, ok := row[dateIdx].(string); ok {
				dateStr = strings.Split(date, "T")[0] // Remove time component
				keyParts = append(keyParts, dateStr)
			}
		}

		// Add dimension value if grouping is specified
		if params.GroupBy != nil && params.GroupBy.Name != nil {
			if idx, ok := columnMap[*params.GroupBy.Name]; ok && len(row) > idx && row[idx] != nil {
				if dimValue, ok := row[idx].(string); ok {
					keyParts = append(keyParts, dimValue)
				}
			}
		}

		// Add second dimension value if second grouping is specified
		if params.GroupBy2 != nil && params.GroupBy2.Name != nil {
			if idx, ok := columnMap[*params.GroupBy2.Name]; ok && len(row) > idx && row[idx] != nil {
				if dimValue, ok := row[idx].(string); ok {
					keyParts = append(keyParts, dimValue)
				}
			}
		}

		if len(keyParts) == 0 {
			continue
		}

		// Build composite key
		rowKey := keyParts[0]
		if len(keyParts) > 1 {
			rowKey += "_" + keyParts[1]
		}
		if len(keyParts) > 2 {
			rowKey += "_" + keyParts[2]
		}

		// Get or create row
		costRow, exists := rowMap[rowKey]
		if !exists {
			costRow = &CostManagementRow{
				Estimated: to.Ptr(false),
			}
			rowMap[rowKey] = costRow

			// Set period dates based on the actual date from the API
			if len(keyParts) > 0 {
				dateStr := keyParts[0]
				costRow.PeriodStart = &dateStr

				// Determine if this is monthly or daily data for period_end calculation
				isMonthlyData := false

				// Check if we found BillingMonth column or if the date is first of month
				if dateColumnName == "BillingMonth" {
					isMonthlyData = true
				} else if len(dateStr) >= 10 && dateStr[8:10] == "01" {
					// Check if the date is the first of the month (likely monthly data)
					isMonthlyData = true
				}

				if isMonthlyData {
					// This appears to be monthly data - set period_end to last day of the month
					if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
						// Calculate last day of the month
						year, month, _ := parsedDate.Date()
						// First day of next month
						nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
						// Last day of current month
						lastDay := nextMonth.AddDate(0, 0, -1)
						lastDayStr := lastDay.Format("2006-01-02")
						costRow.PeriodEnd = &lastDayStr
					} else {
						// Fallback: if parsing fails, use same date
						costRow.PeriodEnd = &dateStr
					}
				} else {
					// This appears to be daily data - period_start and period_end are the same
					costRow.PeriodEnd = &dateStr
				}
			}

			// Set dimension values from keyParts (skip the date which is keyParts[0])
			if len(keyParts) > 1 {
				costRow.Dimension1 = &keyParts[1]
			}
			if len(keyParts) > 2 {
				costRow.Dimension2 = &keyParts[2]
			}
		}

		// Map the available cost data to all expected columns
		// Since we're fetching PreTaxCost, we'll populate multiple columns with this data
		if idx, ok := columnMap["PreTaxCost"]; ok && len(row) > idx && row[idx] != nil {
			if cost, ok := row[idx].(float64); ok {
				// Populate all cost columns with the available PreTaxCost data
				costRow.PreTaxCostAmount = &cost
				costRow.PreTaxCostUnit = to.Ptr("USD")

				// Also populate unblended cost (which is the primary cost metric)
				costRow.UnblendedCostAmount = &cost
				costRow.UnblendedCostUnit = to.Ptr("USD")

				// For monthly data, also populate amortized cost
				// (will be the same as actual cost unless we specifically query amortized data)
				costRow.AmortizedCostAmount = &cost
				costRow.AmortizedCostUnit = to.Ptr("USD")
			}
		}

		// Handle UsageQuantity data if available
		if idx, ok := columnMap["UsageQuantity"]; ok && len(row) > idx && row[idx] != nil {
			if usage, ok := row[idx].(float64); ok {
				costRow.UsageQuantityAmount = &usage
				costRow.UsageQuantityUnit = to.Ptr("Units")
			}
		}

		// Extract currency if available
		if idx, ok := columnMap["Currency"]; ok && len(row) > idx && row[idx] != nil {
			if currency, ok := row[idx].(string); ok {
				costRow.Currency = &currency
				// Update all currency fields to match the actual currency
				if currency != "" {
					costRow.PreTaxCostUnit = &currency
					costRow.UnblendedCostUnit = &currency
					costRow.AmortizedCostUnit = &currency
				}
			}
		}
	}
}
