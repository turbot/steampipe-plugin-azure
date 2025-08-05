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

	// Period dates (optional, populated from query parameters)
	PeriodStart *time.Time
	PeriodEnd   *time.Time

	// Dimension values (populated based on GroupBy)
	Dimension1 *string // Generic dimension field (could be ResourceGroup, ServiceName, etc.)
	Dimension2 *string // Second dimension field for multi-dimensional grouping

	// Cost metrics (following AWS naming conventions)
	// Actual costs (unblended_cost_amount and unblended_cost_unit removed)
	PreTaxCostAmount *float64 // Pre-tax cost (Actual Cost)
	PreTaxCostUnit   *string

	// Amortized costs (for reservations)
	AmortizedCostAmount *float64
	AmortizedCostUnit   *string

	// Metadata
	Estimated *bool
	Currency  *string

	// Common properties
	SubscriptionID   *string
	SubscriptionName *string
	Scope            *string // Azure scope for cost queries
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

// getCostTypeFromString converts cost type string to Azure ExportType
func getCostTypeFromString(costType string) armcostmanagement.ExportType {
	switch costType {
	case "AmortizedCost":
		return armcostmanagement.ExportTypeAmortizedCost
	case "ActualCost":
		return armcostmanagement.ExportTypeActualCost
	default:
		// Default to ActualCost
		return armcostmanagement.ExportTypeActualCost
	}
}

// AllCostMetrics returns all available cost metrics for Azure Cost Management (like AWS)
func AllCostMetrics() []string {
	return []string{
		"PreTaxCost", // Actual Cost
		"Cost",       // Amortized Cost
	}
}

// getMetricsByQueryContext dynamically determines which metrics to fetch based on query columns
func getMetricsByQueryContext(qc *plugin.QueryContext) []string {
	queryColumns := qc.Columns
	var metrics []string

	// Check for cost metrics only (usage quantity removed)
	needsCost := false

	for _, c := range queryColumns {
		switch c {
		case "pre_tax_cost_amount", "pre_tax_cost_unit", "amortized_cost_amount", "amortized_cost_unit":
			needsCost = true
		}
	}

	// Add metrics based on priority (Azure limit: max 2)
	if needsCost {
		metrics = append(metrics, "PreTaxCost")
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

	// Check what the user is querying for (usage quantity removed)
	needsCost := false

	for _, c := range queryColumns {
		switch c {
		case "pre_tax_cost_amount", "pre_tax_cost_unit", "amortized_cost_amount", "amortized_cost_unit":
			needsCost = true
		}
	}

	// Add the primary cost column (PreTaxCost works for all cost types with fallback)
	if needsCost {
		columns = append(columns, to.Ptr("PreTaxCost"))
	}

	// Default to PreTaxCost if no specific columns requested
	if !needsCost {
		columns = append(columns, to.Ptr("PreTaxCost"))
	}

	return columns
}

// getPeriodTimeRange extracts start and end dates from period_start and period_end quals (following AWS pattern)
func getPeriodTimeRange(keyQuals *plugin.QueryData, granularity string) (string, string) {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}

	st, et := "", ""

	// Extract period_start - support multiple operators
	if keyQuals.Quals["period_start"] != nil && len(keyQuals.Quals["period_start"].Quals) <= 1 {
		for _, q := range keyQuals.Quals["period_start"].Quals {
			t := q.Value.GetTimestampValue().AsTime()
			timeStr := t.Format(timeFormat)
			switch q.Operator {
			case "=":
				st = timeStr
			}
		}
	}

	// The API supports a single value with the '=' operator.
	// For queries like: "period_end BETWEEN current_timestamp - interval '31d' AND current_timestamp - interval '1d'", the FDW parses the query parameters with multiple qualifiers.
	// In this case, we will have multiple qualifiers with operators such as:
	// 1. The length of keyQuals.Quals["period_end"].Quals will be 2.
	// 2. The qualifier values would be "2024-05-10" with the '>=" operator and "2024-06-09" with the '<=' operator.
	// In this scenario, manipulating the start and end time is a bit difficult and challenging.
	// Let the API fetch all the rows, and filtering will occur at the Steampipe level.

	// Extract period_end - support multiple operators
	if keyQuals.Quals["period_end"] != nil && len(keyQuals.Quals["period_end"].Quals) <= 1 {
		for _, q := range keyQuals.Quals["period_end"].Quals {
			t := q.Value.GetTimestampValue().AsTime()
			timeStr := t.Format(timeFormat)
			switch q.Operator {
			case "=":
				et = timeStr
			}
		}
	}

	now := time.Now()

	// Set defaults if not provided
	if st == "" {
		st = now.AddDate(0, -11, -30).Format(timeFormat) // 11 months 30 days ago
	}
	if et == "" {
		et = now.AddDate(0, 0, -1).Format(timeFormat) // Yesterday
	}

	return st, et
}

// buildFilterExpression creates filter expressions from table key column quals
func buildFilterExpression(d *plugin.QueryData, dimensionName string) *armcostmanagement.QueryFilter {
	var filters []*armcostmanagement.QueryFilter
	// Process dimension-specific quals (like service_name = 'Storage')
	for _, keyQual := range d.Table.List.KeyColumns {
		if keyQual.Name == "usage_date" || keyQual.Name == "scope" || keyQual.Name == "type" || keyQual.Name == "period_start" || keyQual.Name == "period_end" {
			continue // Skip time, scope, type quals
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
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("UsageDate"),
		},
		{
			Name:        "period_start",
			Description: "The start date of the period, populated if specified in query parameters.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromQual("period_start"),
		},
		{
			Name:        "period_end",
			Description: "The end date of the period, populated if specified in query parameters.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromQual("period_end"),
		},
		{
			Name:        "estimated",
			Description: "Whether the cost data is estimated.",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Estimated"),
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
		{
			Name:        "scope",
			Description: "The Azure scope for the cost query (e.g., subscription, resource group, etc.).",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Scope"),
		},
		{
			Name:        "type",
			Description: "The cost type for the query. Valid values are 'ActualCost' and 'AmortizedCost'. Defaults to 'ActualCost'.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromQual("type"),
			Default:     "ActualCost",
		},
	}

	// Prepend table-specific columns to standard columns
	return append(columns, standardColumns...)
}

// costManagementKeyColumns returns the standard key columns for cost management tables (like AWS)
func costManagementKeyColumns() plugin.KeyColumnSlice {
	return plugin.KeyColumnSlice{
		{
			Name:      "scope",
			Require:   plugin.Optional,
			Operators: []string{"="},
		},
		{
			Name:      "type",
			Require:   plugin.Optional,
			Operators: []string{"="},
		},
		{
			Name:       "period_start",
			Require:    plugin.Optional,
			Operators:  []string{"="},
			CacheMatch: query_cache.CacheMatchExact,
		},
		{
			Name:       "period_end",
			Require:    plugin.Optional,
			Operators:  []string{"="},
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
func streamCostAndUsage(ctx context.Context, d *plugin.QueryData, queryDef armcostmanagement.QueryDefinition, scope string, groupingNames ...string) (interface{}, error) {
	client, err := getCostManagementClient(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	// Resolve scope if it's a placeholder or extract from quals
	if scope == "/subscriptions/placeholder" {
		// Check if scope is provided in quals
		scopeQual := d.EqualsQualString("scope")
		if scopeQual != "" {
			scope = scopeQual
		} else {
			// Fallback to subscription ID if no scope provided
			subscriptionData, err := getSubscriptionID(ctx, d, nil)
			if err != nil {
				return nil, err
			}
			subscriptionID := subscriptionData.(string)
			scope = "/subscriptions/" + subscriptionID
		}
	}

	// Execute query with a single API call
	rowMap := make(map[string]*CostManagementRow)

	plugin.Logger(ctx).Debug("Making Azure Cost Management API call", "query", queryDef, "scope", scope)

	result, err := client.Usage(ctx, scope, queryDef, nil)
	if err != nil {
		plugin.Logger(ctx).Error("Azure Cost Management Query failed", "error", err)
		return nil, err
	}

	processQueryResults(&result.QueryResult, scope, rowMap, groupingNames...)

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
func processQueryResults(result *armcostmanagement.QueryResult, scope string, rowMap map[string]*CostManagementRow, groupingNames ...string) {
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

		// Look for date columns - try UsageDate first, then BillingMonth
		if idx, ok := columnMap["UsageDate"]; ok {
			dateColumnName = "UsageDate"
			dateIdx = idx
		} else if idx, ok := columnMap["BillingMonth"]; ok {
			dateColumnName = "BillingMonth"
			dateIdx = idx
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

		// Add dimension values based on the specific groupings requested
		// This ensures dimensions map correctly to Dimension1 and Dimension2 in the right order
		var dim1Value, dim2Value string

		// Get dimension values in the correct order based on groupingNames
		if len(groupingNames) > 0 {
			// First grouping -> Dimension1
			if idx, ok := columnMap[groupingNames[0]]; ok && len(row) > idx && row[idx] != nil {
				if dimValue, ok := row[idx].(string); ok {
					dim1Value = dimValue
					keyParts = append(keyParts, dimValue)
				}
			}
		}

		if len(groupingNames) > 1 {
			// Second grouping -> Dimension2
			if idx, ok := columnMap[groupingNames[1]]; ok && len(row) > idx && row[idx] != nil {
				if dimValue, ok := row[idx].(string); ok {
					dim2Value = dimValue
					keyParts = append(keyParts, dimValue)
				}
			}
		}

		// If no grouping names provided, fall back to any available dimension columns
		if len(groupingNames) == 0 {
			for colName, idx := range columnMap {
				if colName != "UsageDate" && colName != "BillingMonth" && colName != "Currency" &&
					colName != "PreTaxCost" && colName != "Cost" && colName != "CostUSD" &&
					len(row) > idx && row[idx] != nil {
					if dimValue, ok := row[idx].(string); ok {
						keyParts = append(keyParts, dimValue)
					}
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

				// Set dimension values using the explicitly mapped values
				if dim1Value != "" {
					costRow.Dimension1 = &dim1Value
				}
				if dim2Value != "" {
					costRow.Dimension2 = &dim2Value
				}
			} else {
				// Fallback case: no date column found (should be rare)
				// Use current date as fallback
				fallbackDate := time.Now()
				costRow.UsageDate = &fallbackDate

				// Use the explicitly mapped dimension values
				if dim1Value != "" {
					costRow.Dimension1 = &dim1Value
				}
				if dim2Value != "" {
					costRow.Dimension2 = &dim2Value
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

				// For environments without reservations, use PreTaxCost for amortized costs as fallback
				if costRow.AmortizedCostAmount == nil {
					costRow.AmortizedCostAmount = &cost
					costRow.AmortizedCostUnit = &currency
				}
			}
		}

		// Handle Cost metric (if available from AmortizedCost query)
		// Note: UnblendedCostAmount has been removed from the schema
		if idx, ok := columnMap["Cost"]; ok && len(row) > idx && row[idx] != nil {
			// Cost metric is no longer mapped to unblended_cost_amount as it's been removed
			_ = idx // Suppress unused variable warning
		}

		// Handle CostUSD metric (if available from AmortizedCost query)
		if idx, ok := columnMap["CostUSD"]; ok && len(row) > idx && row[idx] != nil {
			if cost, ok := row[idx].(float64); ok {
				costRow.AmortizedCostAmount = &cost
				costRow.AmortizedCostUnit = &currency
			}
		}

		// Set the currency for all cost fields
		costRow.Currency = &currency

		// Set scope from params
		if costRow.Scope == nil {
			costRow.Scope = &scope
		}
	}
}
