package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2019-11-01/costmanagement"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureCostManagement(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cost_management",
		Description: "Azure Cost Management",
		List: &plugin.ListConfig{
			Hydrate: listCostManagement,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_link",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("QueryProperties.NextLink"),
			},
			{
				Name:        "columns",
				Description: "",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("QueryProperties.Columns"),
			},
			{
				Name:        "rows",
				Description: "",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("QueryProperties.Rows"),
			},
			{
				Name:        "etag",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "location",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// LIST FUNCTION

func listCostManagement(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	costmanagementClient := costmanagement.NewQueryClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	costmanagementClient.Authorizer = session.Authorizer

	// https://learn.microsoft.com/en-us/rest/api/cost-management/query/usage?tabs=HTTP#querycolumn
	// The API needs 2 things - scope (required) and parameters (required)

	/* Parameters: scope - the scope associated with query and export operations. This includes '/subscriptions/{subscriptionId}/' for subscription scope, '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}' for resourceGroup scope, '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}' for Billing Account scope and '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/departments/{departmentId}' for Department scope, '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/enrollmentAccounts/{enrollmentAccountId}' for EnrollmentAccount scope, '/providers/Microsoft.Management/managementGroups/{managementGroupId} for Management Group scope, '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}' for billingProfile scope, '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}/invoiceSections/{invoiceSectionId}' for invoiceSection scope, and '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/customers/{customerId}' specific for partners. parameters - parameters supplied to the CreateOrUpdate Query Config operation. */

	scope := "/subscriptions/" + subscriptionID
	parameters := costmanagement.QueryDefinition{}

	// Type - The type of the query. Possible values include: 'ExportTypeUsage', 'ExportTypeActualCost', 'ExportTypeAmortizedCost'

	// we can keep it as required column (required)
	parameters.Type = "Usage"

	// Timeframe - The time frame for pulling data for the query. If custom, then a specific time period must be provided. Possible values include: 'TimeframeTypeMonthToDate', 'TimeframeTypeBillingMonthToDate', 'TimeframeTypeTheLastMonth', 'TimeframeTypeTheLastBillingMonth', 'TimeframeTypeWeekToDate', 'TimeframeTypeCustom'

	// we can keep it as required column (required)
	parameters.Timeframe = "TimeframeTypeTheLastMonth"

	// parameters.Dataset = &costmanagement.QueryDataset{}
	// at least either Configuration or (Aggregation & Grouping) need to set to get the desired
	/* type QueryDataset struct {
	    // Granularity - The granularity of rows in the query. Possible values include: 'GranularityTypeDaily'
	    Granularity GranularityType `json:"granularity,omitempty"`
	    // Configuration - Has configuration information for the data in the export. The configuration will be ignored if aggregation and grouping are provided.
	    Configuration *QueryDatasetConfiguration `json:"configuration,omitempty"`
	    // Aggregation - Dictionary of aggregation expression to use in the query. The key of each item in the dictionary is the alias for the aggregated column. Query can have up to 2 aggregation clauses.
	    Aggregation map[string]*QueryAggregation `json:"aggregation"`
	    // Grouping - Array of group by expression to use in the query. Query can have up to 2 group by clauses.
	    Grouping *[]QueryGrouping `json:"grouping,omitempty"`
	    // Filter - The filter expression to use in the query. Please reference our Query API REST documentation for how to properly format the filter.
	    Filter *QueryFilter `json:"filter,omitempty"`
	} */

	aggregation := make(map[string]*costmanagement.QueryAggregation)
	agg := &costmanagement.QueryAggregation{
		// valid values are - 'UsageQuantity','PreTaxCost','Cost','CostUSD','PreTaxCostUSD'
		Name: aws.String("PreTaxCost"),

		// Function - The name of the aggregation function to use. Possible values include: 'FunctionTypeAvg', 'FunctionTypeMax', 'FunctionTypeMin', 'FunctionTypeSum'
		Function: "Sum",
	}
	aggregation["total"] = agg // can keep anything here for the string value, it is just a reference name

	grouping := costmanagement.QueryGrouping{
		/*valid values for name: 'ResourceGroup','ResourceGroupName','ResourceLocation','ConsumedService','ResourceType','ResourceId','MeterId','BillingMonth','MeterCategory','MeterSubcategory','Meter','AccountName','DepartmentName','SubscriptionId','SubscriptionName','ServiceName','ServiceTier','EnrollmentAccountName','BillingAccountId','ResourceGuid','BillingPeriod','InvoiceNumber','ChargeType','PublisherType','ReservationId','ReservationName','Frequency','PartNumber','CostAllocationRuleName','MarkupRuleName','PricingModel','BenefitId','BenefitName',' */
		Name: aws.String("ResourceGroup"),

		// Type - Has type of the column to group. Possible values include: 'QueryColumnTypeTag', 'QueryColumnTypeDimension'
		Type: "Dimension",
	}

	// With Aggregation and Grouping the API is returning all the requested columns (PreTaxCost, ResourceGroup)
	// the default columns are UsageDate and Currency
	parameters.Dataset = &costmanagement.QueryDataset{
		// Granularity - The granularity of rows in the query. Possible values include: 'GranularityTypeDaily'
		Granularity: "Daily",
		Aggregation: aggregation,
		Grouping:    &[]costmanagement.QueryGrouping{grouping},
	}

	// filter := &costmanagement.QueryFilter{} is optional, to further optimize the result
	// Here is a sample expression of filter parameter - https://learn.microsoft.com/en-us/rest/api/cost-management/query/usage?tabs=HTTP#subscriptionquery-legacy

	// with configuration - The API does not provide all the requested columns
	//	parameters.Dataset = &costmanagement.QueryDataset{
	// Granularity - The granularity of rows in the query. Possible values include: 'GranularityTypeDaily'
	//	Granularity: "Daily",

	/*	valid values for Configuration columns: 'SubscriptionGuid','ResourceGroup','ResourceLocation','UsageDateTime','MeterCategory','MeterSubcategory','MeterId','MeterName','MeterRegion','UsageQuantity','ResourceRate','PreTaxCost','ConsumedService','ResourceType','InstanceId','TagsDy','OfferId','AdditionalInfo','ServiceInfo1','ServiceInfo2','ServiceName','ServiceTier','Currency','UnitOfMeasure' */
	// 	Configuration: &costmanagement.QueryDatasetConfiguration{
	// 		Columns: &[]string{"PreTaxCost", "ResourceGroup", "UsageDateTime"},
	// 	},
	// }

	/*	> select jsonb_pretty(columns), jsonb_pretty(rows), subscription_id from azure_cost_management
		+------------------------------+-------------------+--------------------------------------+
		| jsonb_pretty                 | jsonb_pretty      | subscription_id                      |
		+------------------------------+-------------------+--------------------------------------+
		| [                            | [                 | d46d723-f95f-4771-bbb5-527234324659c |
		|     {                        |     [             |                                      |
		|         "name": "UsageDate", |         20230611, |                                      |
		|         "type": "Number"     |         "USD"     |                                      |
		|     },                       |     ],            |                                      |
		|     {                        |     [             |                                      |
		|         "name": "Currency",  |         20230612, |                                      |
		|         "type": "String"     |         "USD"     |                                      |
		|     }                        |     ],            |                                      |
		| ]                            |     [             |                                      |
		|                              |         20230613, |                                      |
		|                              |         "USD"     |                                      |
		|                              |     ],            |                                      |
		|                              |     [             |                                      |
		|                              |         20230614, |                                      |
		|                              |         "USD"     |                                      |
		|                              |     ]             |                                      |
		|                              | ]                 |                                      |
		+------------------------------+-------------------+--------------------------------------+
	*/

	result, err := costmanagementClient.Usage(ctx, scope, parameters)
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, result)

	return nil, err
}

// The below is the structure of the columns and rows, need to figure out how we can associate column names with values present in rows.
// with the below data types structure looks like it is not feasible

/*
type QueryProperties struct {
	// NextLink - The link (url) to the next page of results.
	NextLink *string `json:"nextLink,omitempty"`
	// Columns - Array of columns
	Columns *[]QueryColumn `json:"columns,omitempty"`
	// Rows - Array of rows
	Rows *[][]interface{} `json:"rows,omitempty"`
} */
