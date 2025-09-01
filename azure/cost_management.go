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

var CostMetrics = []string{"PreTaxCost"}

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

	// Cost metrics (following AWS naming conventions)
	// Actual costs (unblended_cost_amount and unblended_cost_unit removed)
	PreTaxCostAmount *float64 // Pre-tax cost (Actual Cost)
	PreTaxCostUnit   *string

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

// AllCostMetrics returns all available cost metrics for Azure Cost Management (like AWS)
func AllCostMetrics() []string {
	return append([]string{}, CostMetrics...)
}

// getMetricsByQueryContext dynamically determines which metrics to fetch based on query columns
func getMetricsByQueryContext(qc *plugin.QueryContext) []string {
	// Currently we only request the metrics defined in CostMetrics (e.g., PreTaxCost)
	return append([]string{}, CostMetrics...)
}

// getColumnsFromQueryContext determines which columns to request from Azure API
func getColumnsFromQueryContext(qc *plugin.QueryContext) []*string {
	var columns []*string

	// Always include currency and the metric columns referenced in CostMetrics
	columns = append(columns, to.Ptr("Currency"))
	for _, m := range CostMetrics {
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
		// {
		// 	Name:        "estimated",
		// 	Description: "Whether the cost data is estimated.",
		// 	Type:        proto.ColumnType_BOOL,
		// 	Transform:   transform.FromField("Estimated"),
		// },
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
		// {
		// 	Name:        "amortized_cost_amount",
		// 	Description: "This cost metric reflects the effective cost of the upfront and monthly reservation fees spread across the billing period. By default, Cost Explorer shows the fees for Reserved Instances as a spike on the day that you're charged, but if you choose to show costs as amortized costs, the costs are amortized over the billing period. This means that the costs are broken out into the effective daily rate. Azure estimates your amortized costs by combining your unblended costs with the amortized portion of your upfront and recurring reservation fees.",
		// 	Type:        proto.ColumnType_DOUBLE,
		// 	Transform:   transform.FromField("AmortizedCostAmount"),
		// },
		// {
		// 	Name:        "amortized_cost_unit",
		// 	Description: "Unit type for amortized costs.",
		// 	Type:        proto.ColumnType_STRING,
		// 	Transform:   transform.FromField("AmortizedCostUnit"),
		// },
		{
			Name:        "scope",
			Description: "The Azure scope for the cost query (e.g., subscription, resource group, etc.).",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Scope"),
		},
		{
			Name:        "type",
			Description: "The cost type for the query. Valid values are 'ActualCost' and 'AmortizedCost'.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromQual("type"),
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
		if v, ok := getFloat(row, "PreTaxCost"); ok {
			costRow.PreTaxCostAmount = &v
			costRow.PreTaxCostUnit = &currency
		}
		costRow.Currency = &currency

		// 5) scope
		if costRow.Scope == nil {
			costRow.Scope = &scope
		}
	}
}
