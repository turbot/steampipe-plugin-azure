package azure

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureConsuptionUsage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_consumption_usage",
		Description: "Azure Consuption Usage",
		List: &plugin.ListConfig{
			Hydrate: listConsuptionUsage,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "filter",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
				{
					Name:      "metric",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
				{
					Name:      "scope",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
				{
					Name:      "expand",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The ID that uniquely identifies an event.",
			},
			{
				Name:        "id",
				Description: "The full qualified ARM ID of an event.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "scope",
				Description: "The scope associated with usage details operations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric",
				Description: "Allows to select different type of cost/usage records. Possible values are 'actualcost', 'amortizedcost' or 'usage'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("metric"),
			},
			{
				Name:        "filter",
				Description: "May be used to filter usageDetails by properties/resourceGroup, properties/instanceName, properties/resourceId, properties/chargeType, properties/reservationId, properties/publisherType or tags. The filter supports 'eq', 'lt', 'gt', 'le', 'ge', and 'and'. It does not currently support 'ne', 'or', or 'not'. Tag filter is a key value pair string where key and value is separated by a colon (:). PublisherType Filter accepts two values azure and marketplace and it is currently supported for Web Direct Offer Type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
			},
			{
				Name:        "expand",
				Description: "May be used to expand the 'properties/additionalInfo' or 'properties/meterDetails' within a list of usage details. By default, these fields are not included when listing usage details.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("expand"),
			},
			{
				Name:        "kind",
				Description: "Specifies the kind of usage details.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The etag for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "modern_usage_detail",
				Description: "The modern usage detail.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "legacy_usage_detail",
				Description: "The legacy usage detail.",
				Type:        proto.ColumnType_JSON,
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
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},
		}),
	}
}

type UsageDetails struct {
	Scope             *string
	Kind              consumption.Kind
	ID                *string
	Name              *string
	Type              *string
	Etag              *string
	Tags              map[string]*string
	ModernUsageDetail map[string]interface{}
	LegacyUsageDetail map[string]interface{}
}

//// LIST FUNCTION

func listConsuptionUsage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_consumption_usage.listConsuptionUsage", "sessin_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID

	consumptionClient := consumption.NewUsageDetailsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	consumptionClient.Authorizer = session.Authorizer

	scope := "/subscriptions/" + subscriptionID + "/" // Default scope is subscription
	if d.EqualsQualString("scope") != "" {
		scope = d.EqualsQualString("scope")
	}
	expand := ""
	if d.EqualsQualString("expand") != "" {
		expand = d.EqualsQualString("expand")
	}

	/**
	• The API returns an error if a billing time period is not specified for consumption usage.
	• Error: Billing Period is not supported in (2021-10-01) API Version for Subscription Scope With Web Direct Offer.  Please provide the UsageStart and UsageEnd dates in the $filter key as parameters.
	• For consumption queries, specifying the start and end dates is required.
	• By default, the time period considered is the past year.
	**/

	filter := getConsumptionFilter(d.Quals)
	skiptoken := ""
	var metric consumption.Metrictype
	if d.EqualsQualString("metric") != "" {
		switch d.EqualsQualString("metric") {
		case "actualcost":
			metric = consumption.MetrictypeActualCostMetricType
		case "amortizedcost":
			metric = consumption.MetrictypeAmortizedCostMetricType
		case "usage":
			metric = consumption.MetrictypeUsageMetricType
		}
	}
	result, err := consumptionClient.List(ctx, scope, expand, filter, skiptoken, nil, metric)
	if err != nil {
		plugin.Logger(ctx).Error("azure_consumption_usage.listConsuptionUsage", "api_error", err)
		return nil, err
	}

	for _, res := range result.Values() {
		result := getUsageDetailsByUsageDetailKind(res, scope)
		if result != nil {
			d.StreamListItem(ctx, result)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_consumption_usage.listConsuptionUsage", "paging_error", err)
			return nil, err
		}

		for _, res := range result.Values() {
			result := getUsageDetailsByUsageDetailKind(res, scope)
			if result != nil {
				d.StreamListItem(ctx, result)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, err
}

//// UTILITY FUNCTION

// Get usage details for all type(madern, legacy, basic) of consumption usage.
func getUsageDetailsByUsageDetailKind(res consumption.BasicUsageDetail, scope string) *UsageDetails {
	result := &UsageDetails{}
	modernUsageDetails, isModern := res.AsModernUsageDetail()
	legacyUsageDetails, isLegacy := res.AsLegacyUsageDetail()
	usageDetails, ud := res.AsUsageDetail()
	if isModern {
		result.Scope = &scope
		result.ID = modernUsageDetails.ID
		result.Etag = modernUsageDetails.Etag
		result.Name = modernUsageDetails.Name
		result.Tags = modernUsageDetails.Tags
		result.Type = modernUsageDetails.Type
		result.Kind = modernUsageDetails.Kind
		result.ModernUsageDetail = extractUsageDetailProperties(modernUsageDetails.ModernUsageDetailProperties)
	}
	if isLegacy {
		result.Scope = &scope
		result.ID = legacyUsageDetails.ID
		result.Etag = legacyUsageDetails.Etag
		result.Name = legacyUsageDetails.Name
		result.Tags = legacyUsageDetails.Tags
		result.Type = legacyUsageDetails.Type
		result.Kind = legacyUsageDetails.Kind
		result.LegacyUsageDetail = extractUsageDetailProperties(legacyUsageDetails.LegacyUsageDetailProperties)
	}
	if ud {
		result.Scope = &scope
		result.ID = usageDetails.ID
		result.Etag = usageDetails.Etag
		result.Name = usageDetails.Name
		result.Tags = usageDetails.Tags
		result.Type = usageDetails.Type
		result.Kind = usageDetails.Kind
	}

	return result
}


// When directly accessing an inner attribute using the "FromField()" function, the value is being populated as null even though the response contains a value.
// Therefore, it's necessary to extract the value of the nested attributes from the response.
func extractUsageDetailProperties(value interface{}) map[string]interface{} {
	data := make(map[string]interface{})

	switch item := value.(type) {
	case *consumption.ModernUsageDetailProperties:
		if item != nil {

			// Use reflection to iterate over the struct fields
			data = structToMap(reflect.ValueOf(*item))
		}
	case *consumption.LegacyUsageDetailProperties:
		if item != nil {

			// Use reflection to iterate over the struct fields
			data = structToMap(reflect.ValueOf(*item))
		}
	}

	return data
}

func structToMap(val reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Check if field is a struct and not a zero value
		if fieldValue.Kind() == reflect.Struct && !fieldValue.IsZero() {
			result[field.Name] = structToMap(fieldValue)
		} else if !fieldValue.IsZero() {
			result[field.Name] = fieldValue.Interface()
		} else {
			result[field.Name] = nil
		}
	}

	return result
}

// Construct the filter query parameter in accordance with the API's behavior.
func getConsumptionFilter(quals plugin.KeyColumnQualMap) (filter string) {
	filter = ""
	if quals["filter"] != nil {
		for _, q := range quals["filter"].Quals {
			if q.Operator == "=" {
				val := q.Value.GetStringValue()
				filter = val
			}
		}
	}

	outputLayout := "2006-01-02T15:04:05Z"
	endData := time.Now().AddDate(-1, 0, 0).Format(outputLayout)
	startData := time.Now().Format(outputLayout)

	// Default time period is last one year
	if filter != "" && !strings.Contains(filter, "properties/usageEnd") && !strings.Contains(filter, "properties/usageStart") {
		filter = filter + " and properties/usageEnd eq '" + endData + "' and properties/usageStart eq '" + startData + "'"
	}
	if filter == "" && !strings.Contains(filter, "properties/usageEnd") && !strings.Contains(filter, "properties/usageStart") {
		filter = "properties/usageEnd eq '" + endData + "' and properties/usageStart eq '" + startData + "'"
	}

	return filter
}
