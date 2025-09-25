package azure

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

var CostMetrics = []string{"PreTaxCost", "Cost"}

// CostManagementRow represents a flattened cost management result row with all cost types (like AWS)
// https://learn.microsoft.com/en-us/azure/cost-management-billing/automate/understand-usage-details-fields
type CostManagementRow struct {
	UsageDate *time.Time

	// Period dates (optional, populated from query parameters)
	PeriodStart *time.Time
	PeriodEnd   *time.Time

	// Dimension values (populated based on GroupBy)
	Dimension1 *string // Generic dimension field (could be ResourceGroup, ServiceName, etc.)
	Dimension2 *string // Second dimension field for multi-dimensional grouping
	Dimensions *map[string]string

	// Cost metrics
	PreTaxCostAmount *float64
	CostAmount       *float64

	// Metadata
	Currency *string

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

// getMetricsByQueryContext determines which metrics to fetch based on selected columns
func getMetricsByQueryContext(qc *plugin.QueryContext) []string {
	var metrics []string
	// If user selected cost, include Cost
	costMetricAdded := false
	preTaxCostMetricAdded := false
	if slices.Contains(qc.Columns, "cost") {
		costMetricAdded = true
		metrics = append(metrics, "Cost")
	}

	if slices.Contains(qc.Columns, "pre_tax_cost") {
		preTaxCostMetricAdded = true
		metrics = append(metrics, "PreTaxCost")
	}

	// default to PreTaxCost if nothing explicitly requested
	if !costMetricAdded && !preTaxCostMetricAdded {
		metrics = append(metrics, "PreTaxCost")
	}

	return metrics
}

// getColumnsFromQueryContext determines which columns to request from Azure API
// TODO: Is it required?
func getColumnsFromQueryContext(qc *plugin.QueryContext) []*string {
	var columns []*string
	columns = append(columns, to.Ptr("Currency"))
	for _, m := range getMetricsByQueryContext(qc) {
		columns = append(columns, to.Ptr(m))
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
		if keyQual.Name == "usage_date" || keyQual.Name == "scope" || keyQual.Name == "cost_type" || keyQual.Name == "period_start" || keyQual.Name == "period_end" {
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
		// estimated removed
		{
			Name:        "pre_tax_cost",
			Description: "Pre-tax cost amount for the period.",
			Type:        proto.ColumnType_DOUBLE,
			Transform:   transform.FromField("PreTaxCostAmount"),
		},
		{
			Name:        "cost",
			Description: "Aggregated cost amount for the period (Cost metric).",
			Type:        proto.ColumnType_DOUBLE,
			Transform:   transform.FromField("CostAmount"),
		},
		{
			Name:        "currency",
			Description: "Currency code for the returned cost values.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Currency"),
		},
		{
			Name:        "scope",
			Description: "The Azure scope for the cost query (e.g., subscription, resource group, etc.).",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Scope"),
		},
		{
			Name:        "cost_type",
			Description: "The cost type for the query. Valid values are 'ActualCost' and 'AmortizedCost'.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromQual("cost_type"),
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
			Name:      "cost_type",
			Require:   plugin.Required,
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

// getForecastClient creates a new forecast client
func getForecastClient(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (*armcostmanagement.ForecastClient, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	client, err := armcostmanagement.NewForecastClient(session.Cred, session.ClientOptions)
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

	// Build column index map for quick lookups
	columnIdx := make(map[string]int)
	for i, col := range result.Properties.Columns {
		if col.Name != nil {
			columnIdx[*col.Name] = i
		}
	}

	// Helpers
	getString := func(row []any, name string) string {
		if idx, ok := columnIdx[name]; ok && idx < len(row) && row[idx] != nil {
			if s, ok := row[idx].(string); ok {
				return s
			}
		}
		return ""
	}
	getFloat := func(row []any, name string) (float64, bool) {
		if idx, ok := columnIdx[name]; ok && idx < len(row) && row[idx] != nil {
			switch v := row[idx].(type) {
			case float64:
				return v, true
			case int:
				return float64(v), true
			}
		}
		return 0, false
	}

	// Determine date column preference
	dateCol := "UsageDate"
	if _, ok := columnIdx["UsageDate"]; !ok {
		if _, ok2 := columnIdx["BillingMonth"]; ok2 {
			dateCol = "BillingMonth"
		}
	}

	for _, row := range result.Properties.Rows {
		// 1) usage_date
		var usageDateStr string
		if dateCol == "UsageDate" {
			usageDateStr = getString(row, "UsageDate")
			if usageDateStr == "" {
				if n, ok := getFloat(row, "UsageDate"); ok {
					s := fmt.Sprintf("%.0f", n)
					if t, err := time.Parse("20060102", s); err == nil {
						usageDateStr = t.Format("2006-01-02")
					}
				}
			} else if strings.Contains(usageDateStr, "T") {
				usageDateStr = strings.Split(usageDateStr, "T")[0]
			}
		} else { // BillingMonth -> normalize to last day of month
			usageDateStr = getString(row, "BillingMonth")
			if usageDateStr != "" {
				// Try RFC3339 first
				if t, err := time.Parse(time.RFC3339, usageDateStr); err == nil {
					y, m, _ := t.Date()
					next := time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC)
					usageDateStr = next.AddDate(0, 0, -1).Format("2006-01-02")
				} else {
					// Handle formats like 2006-01-02T15:04:05 (no timezone) or 2006-01-02
					if strings.Contains(usageDateStr, "T") {
						usageDateStr = strings.Split(usageDateStr, "T")[0]
					}
					if t2, err2 := time.Parse("2006-01-02", usageDateStr); err2 == nil {
						y, m, _ := t2.Date()
						next := time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC)
						usageDateStr = next.AddDate(0, 0, -1).Format("2006-01-02")
					} else {
						// leave empty and try numeric path below
						usageDateStr = ""
					}
				}
			}
			if usageDateStr == "" { // try numeric encodings
				if n, ok := getFloat(row, "BillingMonth"); ok {
					s := fmt.Sprintf("%.0f", n)
					// try YYYYMM, fallback YYYYMMDD
					if len(s) == 6 {
						if t, err := time.Parse("200601", s); err == nil {
							y, m, _ := t.Date()
							next := time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC)
							usageDateStr = next.AddDate(0, 0, -1).Format("2006-01-02")
						}
					} else if len(s) == 8 {
						if t, err := time.Parse("20060102", s); err == nil {
							y, m, _ := t.Date()
							next := time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC)
							usageDateStr = next.AddDate(0, 0, -1).Format("2006-01-02")
						}
					}
				}
			}
		}
		if usageDateStr == "" {
			continue
		}
		parsedDate, err := time.Parse("2006-01-02", usageDateStr)
		if err != nil {
			continue
		}

		// 2) dimensions in the order provided
		var dim1, dim2 string
		if len(groupingNames) > 0 {
			dim1 = getString(row, groupingNames[0])
		}
		if len(groupingNames) > 1 {
			dim2 = getString(row, groupingNames[1])
		}

		// 3) stable key
		rowKey := usageDateStr + "|" + dim1 + "|" + dim2
		costRow, exists := rowMap[rowKey]
		if !exists {
			costRow = &CostManagementRow{}
			costRow.UsageDate = &parsedDate
			if dim1 != "" {
				costRow.Dimension1 = &dim1
			}
			if dim2 != "" {
				costRow.Dimension2 = &dim2
			}
			rowMap[rowKey] = costRow
		}

		if len(groupingNames) > 0 {
			costRow.Dimensions = &map[string]string{}
			for _, dimName := range groupingNames {
				(*costRow.Dimensions)[dimName] = getString(row, dimName)
			}
		}

		// 4) currency and cost mapping
		currency := getString(row, "Currency")
		if currency == "" {
			currency = "USD"
		}
		// Handle PreTaxCost metric
		if idx, ok := columnIdx["PreTaxCost"]; ok && len(row) > idx && row[idx] != nil {
			if cost, ok := row[idx].(float64); ok {
				costRow.PreTaxCostAmount = &cost
			}
		}

		// Handle Cost metric
		if idx, ok := columnIdx["Cost"]; ok && len(row) > idx && row[idx] != nil {
			if v, ok := row[idx].(float64); ok {
				costRow.CostAmount = &v
			}
		}

		// Set the currency for all cost fields
		costRow.Currency = &currency

		// 5) scope
		if costRow.Scope == nil {
			costRow.Scope = &scope
		}
	}
}

// buildForecastQueryInput builds input parameters specifically for forecast tables
func buildForecastQueryInput(ctx context.Context, d *plugin.QueryData, granularity string) (armcostmanagement.ForecastDefinition, string, error) {
	// Get scope from quals
	scope := d.EqualsQualString("scope")
	if scope == "" {
		// Get subscription ID
		subscriptionID, err := getSubscriptionID(ctx, d, nil)
		if err != nil {
			return armcostmanagement.ForecastDefinition{}, "", err
		}
		scope = "/subscriptions/" + subscriptionID.(string)
	}

	// Get cost type from quals
	costType := d.EqualsQualString("cost_type")
	if costType == "" {
		return armcostmanagement.ForecastDefinition{}, "", fmt.Errorf("missing required qual 'cost_type' (ActualCost | AmortizedCost)")
	}

	// Set timeframe and forecast period
	timeframe := armcostmanagement.TimeframeTypeCustom
	startDate := time.Now().UTC()
	var endDate time.Time

	// Set forecast period based on granularity
	switch granularity {
	case "Monthly":
		endDate = startDate.AddDate(1, 0, 0) // 1 year forecast for monthly
	case "Daily":
		endDate = startDate.AddDate(0, 3, 0) // 3 months forecast for daily
	}

	// Handle user-provided dates
	if d.Quals["period_start"] != nil {
		for _, q := range d.Quals["period_start"].Quals {
			ts := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=", ">=", ">":
				if ts.After(startDate) {
					startDate = ts
				}
			}
		}
	}
	if d.Quals["period_end"] != nil {
		for _, q := range d.Quals["period_end"].Quals {
			ts := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=", "<=", "<":
				if ts.Before(endDate) {
					endDate = ts
				}
			}
		}
	}

	// Ensure we're not forecasting in the past
	if startDate.Before(time.Now().UTC()) {
		startDate = time.Now().UTC()
	}
	if endDate.Before(startDate) {
		switch granularity {
		case "Monthly":
			endDate = startDate.AddDate(1, 0, 0)
		case "Daily":
			endDate = startDate.AddDate(0, 3, 0)
		}
	}

	// Create forecast definition
	forecastDef := armcostmanagement.ForecastDefinition{
		Type:      to.Ptr(armcostmanagement.ForecastType(costType)),
		Timeframe: to.Ptr(armcostmanagement.ForecastTimeframe(timeframe)),
		TimePeriod: &armcostmanagement.ForecastTimePeriod{
			From: to.Ptr(startDate),
			To:   to.Ptr(endDate),
		},
		Dataset: &armcostmanagement.ForecastDataset{
			Granularity: to.Ptr(armcostmanagement.GranularityType(granularity)),
			Aggregation: map[string]*armcostmanagement.ForecastAggregation{
				"PreTaxCost": {
					Name:     to.Ptr(armcostmanagement.FunctionNameCost),
					Function: to.Ptr(armcostmanagement.FunctionTypeSum),
				},
			},
		},
	}

	return forecastDef, scope, nil
}

// streamForecastResults handles forecast API results specifically
func streamForecastResults(ctx context.Context, d *plugin.QueryData, result *armcostmanagement.ForecastClientUsageResponse, scope string, granularity string) error {
	if result.Properties == nil || result.Properties.Rows == nil {
		return nil
	}

	for _, row := range result.Properties.Rows {
		if len(row) < 4 {
			continue
		}

		// Parse the forecast row
		costRow := &CostManagementRow{
			Scope: &scope,
		}

		// Parse date - can be either YYYYMMDD format (daily) or RFC3339/BillingMonth (monthly)
		var usageDate time.Time

		// Try string format first (RFC3339 or BillingMonth)
		if dateStr, ok := row[1].(string); ok {
			// Try RFC3339 first
			if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
				usageDate = t
			} else {
				// Try BillingMonth format (2025-09-01T00:00:00)
				if t, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
					usageDate = t
				}
			}
		} else if dateNum, ok := row[1].(float64); ok {
			// Try YYYYMMDD format
			dateStr := fmt.Sprintf("%.0f", dateNum)
			if t, err := time.Parse("20060102", dateStr); err == nil {
				usageDate = t
			}
		}

		if !usageDate.IsZero() {
			// Set usage date based on granularity
			if strings.HasSuffix(granularity, "Monthly") {
				// For monthly, set to start of month
				startOfMonth := time.Date(usageDate.Year(), usageDate.Month(), 1, 0, 0, 0, 0, usageDate.Location())
				costRow.UsageDate = &startOfMonth
			} else {
				// For daily, set to start of day
				startOfDay := time.Date(usageDate.Year(), usageDate.Month(), usageDate.Day(), 0, 0, 0, 0, usageDate.Location())
				costRow.UsageDate = &startOfDay
			}
		}

		// Parse cost
		if cost, ok := row[0].(float64); ok {
			costRow.PreTaxCostAmount = &cost
			costRow.CostAmount = &cost // For forecast, both values are the same
		}

		// Parse currency
		if currency, ok := row[3].(string); ok {
			costRow.Currency = &currency
		}

		// Extract subscription details from scope
		if strings.HasPrefix(scope, "/subscriptions/") {
			subID := strings.TrimPrefix(scope, "/subscriptions/")
			if idx := strings.Index(subID, "/"); idx != -1 {
				subID = subID[:idx]
			}
			costRow.SubscriptionID = &subID
		}

		d.StreamListItem(ctx, costRow)

		if d.RowsRemaining(ctx) == 0 {
			return nil
		}
	}

	return nil
}

// buildCostQueryInput is a common function to build input parameters for all cost tables
func buildCostQueryInput(ctx context.Context, d *plugin.QueryData, granularity string, groupingNames []string) (armcostmanagement.QueryDefinition, string, error) {
	// Get scope from quals, default to placeholder if not provided
	scope := d.EqualsQualString("scope")
	if scope == "" {
		scope = "/subscriptions/placeholder" // Will be resolved in streamCostAndUsage
	}

	// Get cost type from quals, required
	costType := d.EqualsQualString("cost_type")
	if costType == "" {
		return armcostmanagement.QueryDefinition{}, "", fmt.Errorf("missing required qual 'cost_type' (ActualCost | AmortizedCost)")
	}

	// Set timeframe and time period
	timeframe := armcostmanagement.TimeframeTypeCustom
	timePeriod := &armcostmanagement.QueryTimePeriod{}

	// Get time range from period_start/period_end quals
	startTime, endTime := getPeriodTimeRange(d, granularity)

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

	azureGranularity := getGranularityFromString(granularity)

	// Build GroupBy for specified grouping names
	var groupings []*armcostmanagement.QueryGrouping
	for _, groupName := range groupingNames {
		groupings = append(groupings, &armcostmanagement.QueryGrouping{
			Type: to.Ptr(armcostmanagement.QueryColumnTypeDimension),
			Name: to.Ptr(groupName),
		})
	}

	// Build filter expressions from quals (use first grouping name for filter)
	var filter *armcostmanagement.QueryFilter
	if len(groupingNames) > 0 {
		filter = buildFilterExpression(d, groupingNames[0])
	}

	// Build aggregation based on requested columns
	aggregation := make(map[string]*armcostmanagement.QueryAggregation)

	metrics := getMetricsByQueryContext(d.QueryContext)
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
	}

	// Add grouping if specified
	if len(groupings) > 0 {
		dataset.Grouping = groupings
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

	// Set TimePeriod for custom timeframe
	queryDef.TimePeriod = timePeriod

	return queryDef, scope, nil
}
