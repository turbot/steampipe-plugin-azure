package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITON

func tableAzureResourceResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_resource",
		Description: "Azure Resource",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getResource,
			Tags: map[string]string{
				"service": "Microsoft.Resources",
				"action":  "resources/read",
			},
			// No error is returned if the resource is not found
		},
		List: &plugin.ListConfig{
			Hydrate: listResources,
			Tags: map[string]string{
				"service": "Microsoft.Resources",
				"action":  "resources/read",
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "region", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "type", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "identity_principal_id", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "plan_publisher", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "plan_name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "plan_product", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "plan_promotion_code", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "plan_version", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "resource_group", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "filter", Require: plugin.Optional, Operators: []string{"="}, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "Resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "The filter to apply on the operation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
			},
			{
				Name:        "created_time",
				Description: "The created time of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreatedTime").Transform(transform.NullIfZeroValue).Transform(convertDateToTime),
			},
			{
				Name:        "changed_time",
				Description: "The changed time of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ChangedTime").Transform(transform.NullIfZeroValue).Transform(convertDateToTime),
			},
			{
				Name:        "identity_principal_id",
				Description: "The principal ID of resource identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.PrincipalID"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "plan_publisher",
				Description: "The plan publisher ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Plan.Publisher"),
			},
			{
				Name:        "plan_name",
				Description: "The plan ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Plan.Name"),
			},
			{
				Name:        "plan_product",
				Description: "The plan offer ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Plan.Product"),
			},
			{
				Name:        "plan_promotion_code",
				Description: "The plan promotion code.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Plan.PromotionCode"),
			},
			{
				Name:        "plan_version",
				Description: "The plan's version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Plan.Version"),
			},
			{
				Name:        "kind",
				Description: "The kind of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_by",
				Description: "ID of the resource that manages this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku",
				Description: "The SKU of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "The identity of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "extended_location",
				Description: "Resource extended location.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "properties",
				Description: "The resource properties.",
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

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_resource.listResources", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	resourceClient := resources.NewClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	resourceClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &resourceClient, d.Connection)

	// https://learn.microsoft.com/en-us/rest/api/resources/resources/list?view=rest-resources-2021-04-01#uri-parameters
	filter := getResourceFilter(d.Quals)
	expand := "createdTime,changedTime,provisioningState"
	// A value less than or equal to '1000' must be provided.
	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	result, err := resourceClient.List(ctx, filter, expand, &maxLimit)
	if err != nil {
		plugin.Logger(ctx).Error("azure_resource.listResources", "api_error", err)
		return nil, err
	}
	for _, resource := range result.Values() {
		d.StreamListItem(ctx, resource)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_resource.listResources", "api_paging_error", err)
			return nil, err
		}
		for _, resource := range result.Values() {
			d.StreamListItem(ctx, resource)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_resource.getResource", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	id := d.EqualsQuals["id"].GetStringValue()
	if id == "" {
		return nil, nil
	}

	apiVersion := "2021-04-01"

	resourceClient := resources.NewClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	resourceClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &resourceClient, d.Connection)

	op, err := resourceClient.GetByID(ctx, id, apiVersion)
	if err != nil {
		plugin.Logger(ctx).Error("azure_resource.getResource", "api_error", err)
		return nil, nil
	}

	return op, nil
}

//// UTILITY FUNCTION

// Construct the filter query parameter in accordance with the API's behavior.
func getResourceFilter(quals plugin.KeyColumnQualMap) string {
	if filter := getDirectFilter(quals["filter"]); filter != "" {
		return filter
	}

	filterQuals := map[string]string{
		"region":                "location",
		"type":                  "resourceType",
		"name":                  "name",
		"resource_group":        "resourceGroup",
		"identity_principal_id": "identity/principalId",
		"plan_publisher":        "plan/publisher",
		"plan_name":             "plan/name",
		"plan_product":          "plan/product",
		"plan_promotion_code":   "plan/promotionCode",
		"plan_version":          "plan/version",
	}

	var filters []string
	for columnName, filterValue := range filterQuals {
		if quals[columnName] != nil {
			for _, q := range quals[columnName].Quals {
				if operator, valid := getOperator(q.Operator); valid {
					filters = append(filters, fmt.Sprintf("%s %s '%s'", filterValue, operator, q.Value.GetStringValue()))
				}
			}
		}
	}

	return strings.Join(filters, " and ")
}

// getDirectFilter returns the direct filter value if present
func getDirectFilter(filterQual *plugin.KeyColumnQuals) string {
	if filterQual != nil {
		for _, q := range filterQual.Quals {
			if q.Operator == "=" {
				return q.Value.GetStringValue()
			}
		}
	}
	return ""
}

// getOperator maps the operator to its string representation
func getOperator(operator string) (string, bool) {
	switch operator {
	case "=":
		return "eq", true
	case "<>":
		return "ne", true
	default:
		return "", false
	}
}
