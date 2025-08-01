package azure

import (
	"context"
	"fmt"
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
	UsageDate *time.Time

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

// getMetricsByQueryContext dynamically determines which metrics to fetch based on query columns
func getMetricsByQueryContext(qc *plugin.QueryContext) []string {
	queryColumns := qc.Columns
	var metrics []string

	// Always prioritize PreTaxCost as the primary cost metric (from ActualCost query)
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

	// Add metrics based on priority (Azure limit: max 2)
	if needsCost {
		metrics = append(metrics, "PreTaxCost")
	}
	if needsUsage && len(metrics) < 2 {
		metrics = append(metrics, "UsageQuantity")
	}

	// Default to PreTaxCost if no specific columns requested
	if len(metrics) == 0 {
		metrics = append(metrics, "PreTaxCost")
	}

	return metrics
}

// getColumnsFromQueryContext determines which columns to request from Azure API
func getColumnsFromQueryContext(qc *plugin.QueryContext) []*string {
	queryColumns := qc.Columns
	var columns []*string

	// Always include basic columns
	columns = append(columns, to.Ptr("Currency"))

	// Check what the user is querying for
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

	// Add the primary cost column (PreTaxCost works for all cost types with fallback)
	if needsCost {
		columns = append(columns, to.Ptr("PreTaxCost"))
	}
	if needsUsage {
		columns = append(columns, to.Ptr("UsageQuantity"))
	}

	// Default to PreTaxCost if no specific columns requested
	if !needsCost && !needsUsage {
		columns = append(columns, to.Ptr("PreTaxCost"))
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

// getUsageDateTimeRange extracts start and end dates from usage_date quals (following AWS pattern)
func getUsageDateTimeRange(keyQuals *plugin.QueryData, granularity string) (string, string) {
	timeFormat := "2006-01-02"

	st, et := "", ""

	if keyQuals.Quals["usage_date"] != nil && len(keyQuals.Quals["usage_date"].Quals) <= 1 {
		for _, q := range keyQuals.Quals["usage_date"].Quals {
			t := q.Value.GetTimestampValue().AsTime().Format(timeFormat)
			switch q.Operator {
			case "=", ">=", ">":
				st = t
			case "<", "<=":
				et = t
			}
		}
	}

	now := time.Now()

	// Set defaults if not provided
	if st == "" {
		st = now.AddDate(0, -11, 0).Format(timeFormat) // 11 months ago
	}
	if et == "" {
		et = now.AddDate(0, 0, -1).Format(timeFormat) // Yesterday
	}

	// Parse dates to ensure the range doesn't exceed 1 year
	startDate, err := time.Parse(timeFormat, st)
	if err == nil {
		endDate, err := time.Parse(timeFormat, et)
		if err == nil {
			// If the range exceeds 365 days, adjust the end date
			maxEndDate := startDate.AddDate(1, 0, -1) // 1 year minus 1 day from start
			if endDate.After(maxEndDate) {
				et = maxEndDate.Format(timeFormat)
			}
		}
	}

	return st, et
}

// buildFilterExpression creates filter expressions from table key column quals
func buildFilterExpression(d *plugin.QueryData, dimensionName string) *armcostmanagement.QueryFilter {
	var filters []*armcostmanagement.QueryFilter
	// Process dimension-specific quals (like service_name = 'Storage')
	for _, keyQual := range d.Table.List.KeyColumns {
		if keyQual.Name == "usage_date" || keyQual.Name == "subscription_id" {
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
			Name:        "usage_date",
			Description: "The date for which the cost metric is calculated.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UsageDate"),
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
			Name:       "usage_date",
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

	// Execute query with a single API call, respecting Azure's 2-aggregation limit
	rowMap := make(map[string]*CostManagementRow)

	// Use ActualCost by default for most reliable results
	queryDef := armcostmanagement.QueryDefinition{
		Type:      to.Ptr(armcostmanagement.ExportTypeActualCost),
		Timeframe: to.Ptr(params.Timeframe),
		Dataset:   dataset,
	}

	// Set TimePeriod if using Custom timeframe
	if params.Timeframe == armcostmanagement.TimeframeTypeCustom && params.TimePeriod != nil {
		queryDef.TimePeriod = params.TimePeriod
	}

	plugin.Logger(ctx).Debug("Making Azure Cost Management API call", "query", queryDef)

	result, err := client.Usage(ctx, params.Scope, queryDef, nil)
	if err != nil {
		plugin.Logger(ctx).Error("Azure Cost Management Query failed", "error", err)
		return nil, err
	}

	for _, c := range result.Properties.Columns{

		plugin.Logger(ctx).Error("Query Result Columns ===>>>", *c.Name, *c.Type)
	}
	plugin.Logger(ctx).Error("Query Result Rows ===>>>", result.Properties.Rows)

	processQueryResults(&result.QueryResult, params, "actual", rowMap)

	// Handle pagination
	for result.QueryResult.Properties != nil && result.QueryResult.Properties.NextLink != nil {
		nextURL := *result.QueryResult.Properties.NextLink
		plugin.Logger(ctx).Debug("Following pagination", "nextLink", nextURL)

		result, err = client.Usage(ctx, params.Scope, queryDef, nil)
		if err != nil {
			plugin.Logger(ctx).Error("Pagination failed", "error", err)
			break
		}
		processQueryResults(&result.QueryResult, params, "actual", rowMap)
	}

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

		// Always look for UsageDate column for daily granularity
		if params.Granularity == armcostmanagement.GranularityTypeDaily {
			// Look for UsageDate column (always available in daily data)
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

		// Get date value - only if we have a valid date column
		var dateStr string
		if dateIdx != -1 && len(row) > dateIdx && row[dateIdx] != nil {
			// Handle both string and numeric date formats
			switch date := row[dateIdx].(type) {
			case string:
				// String date format (like "2024-08-01T00:00:00Z")
				dateStr = strings.Split(date, "T")[0] // Remove time component
				keyParts = append(keyParts, dateStr)
			case float64:
				// Numeric date format (like 20240801)
				// Convert to string and parse with Go's time.Parse
				numericDateStr := fmt.Sprintf("%.0f", date)
				if parsedTime, err := time.Parse("20060102", numericDateStr); err == nil {
					dateStr = parsedTime.Format("2006-01-02")
					keyParts = append(keyParts, dateStr)
				}
			case int:
				// Handle int format as well
				numericDateStr := fmt.Sprintf("%d", date)
				if parsedTime, err := time.Parse("20060102", numericDateStr); err == nil {
					dateStr = parsedTime.Format("2006-01-02")
					keyParts = append(keyParts, dateStr)
				}
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

			// Set period dates and dimensions based on what we found
			if dateIdx != -1 && len(keyParts) > 0 {
				// We found a valid date column - standard case for daily and monthly data
				dateStr := keyParts[0]
				if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
					costRow.UsageDate = &parsedDate

					// Determine if this is monthly or daily data for period_end calculation
					isMonthlyData := false
					if dateColumnName == "BillingMonth" {
						isMonthlyData = true
					} else if len(dateStr) >= 10 && dateStr[8:10] == "01" {
						isMonthlyData = true
					}

					if isMonthlyData {
						// Monthly data - set usage_date to last day of the month
						year, month, _ := parsedDate.Date()
						nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
						lastDay := nextMonth.AddDate(0, 0, -1)
						costRow.UsageDate = &lastDay
					}
					// For daily data, we keep the original parsed date
				}

				// Set dimension values (skip the date which is keyParts[0])
				if len(keyParts) > 1 {
					costRow.Dimension1 = &keyParts[1]
				}
				if len(keyParts) > 2 {
					costRow.Dimension2 = &keyParts[2]
				}
			} else {
				// Fallback case: no date column found (should be rare)
				// Use the query time period for date information
				var periodStart string
				if params.TimePeriod != nil && params.TimePeriod.From != nil {
					periodStart = params.TimePeriod.From.Format("2006-01-02")
				} else {
					// Fallback to current date range
					periodStart = time.Now().Format("2006-01-02")
				}
				if parsedPeriodStart, err := time.Parse("2006-01-02", periodStart); err == nil {
					costRow.UsageDate = &parsedPeriodStart
				}

				// ALL keyParts are dimensions in this case (no date in keyParts)
				if len(keyParts) > 0 {
					costRow.Dimension1 = &keyParts[0]
				}
				if len(keyParts) > 1 {
					costRow.Dimension2 = &keyParts[1]
				}
			}
		}

		// Map the different cost types to their appropriate columns based on costType
		currency := "USD" // Default currency

		// Extract currency if available
		if idx, ok := columnMap["Currency"]; ok && len(row) > idx && row[idx] != nil {
			if curr, ok := row[idx].(string); ok && curr != "" {
				currency = curr
			}
		}

		// Handle different cost metrics based on what's available in the API response

		// Handle PreTaxCost metric (primary from ActualCost query)
		if idx, ok := columnMap["PreTaxCost"]; ok && len(row) > idx && row[idx] != nil {
			if cost, ok := row[idx].(float64); ok {
				costRow.PreTaxCostAmount = &cost
				costRow.PreTaxCostUnit = &currency

				// For environments without reservations, use PreTaxCost for all cost types
				if costRow.UnblendedCostAmount == nil {
					costRow.UnblendedCostAmount = &cost
					costRow.UnblendedCostUnit = &currency
				}
				if costRow.AmortizedCostAmount == nil {
					costRow.AmortizedCostAmount = &cost
					costRow.AmortizedCostUnit = &currency
				}
			}
		}

		// Handle Cost metric (if available from AmortizedCost query)
		if idx, ok := columnMap["Cost"]; ok && len(row) > idx && row[idx] != nil {
			if cost, ok := row[idx].(float64); ok {
				costRow.UnblendedCostAmount = &cost
				costRow.UnblendedCostUnit = &currency
			}
		}

		// Handle CostUSD metric (if available from AmortizedCost query)
		if idx, ok := columnMap["CostUSD"]; ok && len(row) > idx && row[idx] != nil {
			if cost, ok := row[idx].(float64); ok {
				costRow.AmortizedCostAmount = &cost
				costRow.AmortizedCostUnit = &currency
			}
		}

		// Handle UsageQuantity metric
		if idx, ok := columnMap["UsageQuantity"]; ok && len(row) > idx && row[idx] != nil {
			if usage, ok := row[idx].(float64); ok {
				costRow.UsageQuantityAmount = &usage
				costRow.UsageQuantityUnit = &currency
			}
		}

		// Set the currency for all cost fields
		costRow.Currency = &currency
	}
}
