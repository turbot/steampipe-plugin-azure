package azure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resourcegraph/mgmt/resourcegraph"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureResourceGraph(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_resource_graph",
		Description: "Execute an Azure Resource Graph query and return the results as rows.",
		List: &plugin.ListConfig{
			Hydrate: listAzureResourceGraph,
			Tags: map[string]string{
				"service": "Microsoft.ResourceGraph",
				"action":  "resources/read",
			},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "query",
					Require:   plugin.Required,
					Operators: []string{"="},
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The resource ID, if projected by the query.",
				Transform:   transform.FromField("id"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The resource name, if projected by the query.",
				Transform:   transform.FromField("name"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The resource type, if projected by the query.",
				Transform:   transform.FromField("type"),
			},
			{
				Name:        "kind",
				Type:        proto.ColumnType_STRING,
				Description: "The kind of the resource, if available.",
				Transform:   transform.FromField("kind"),
			},
			{
				Name:        "identity",
				Type:        proto.ColumnType_JSON,
				Description: "The managed identity info of the resource, if available.",
				Transform:   transform.FromField("identity"),
			},
			{
				Name:        "managed_by",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the resource that manages this resource, if available.",
				Transform:   transform.FromField("managedBy"),
			},
			{
				Name:        "plan",
				Type:        proto.ColumnType_JSON,
				Description: "The plan info of the resource, if available.",
				Transform:   transform.FromField("plan"),
			},
			{
				Name:        "properties",
				Type:        proto.ColumnType_JSON,
				Description: "The resource properties as returned by the graph query.",
				Transform:   transform.FromField("properties"),
			},
			{
				Name:        "sku",
				Type:        proto.ColumnType_JSON,
				Description: "The SKU of the resource, if available.",
				Transform:   transform.FromField("sku"),
			},
			{
				Name:        "tenant_id",
				Type:        proto.ColumnType_STRING,
				Description: "The tenant ID of the resource, if available.",
				Transform:   transform.FromField("tenantId"),
			},

			{
				Name:        "zones",
				Type:        proto.ColumnType_JSON,
				Description: "The availability zones of the resource, if available.",
				Transform:   transform.FromField("zones"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.FromField("tags"),
			},
			{
				Name:        "extended_location",
				Type:        proto.ColumnType_JSON,
				Description: "The extended location info of the resource, if available.",
				Transform:   transform.FromField("extendedLocation"),
			},
			{
				Name:        "query",
				Type:        proto.ColumnType_STRING,
				Description: "The KQL query executed against Azure Resource Graph.",
				Transform:   transform.FromQual("query"),
			},

			// Azure standard columns
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionRegion,
				Transform:   transform.FromField("location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionResourceGroup,
				Transform:   transform.FromField("resourceGroup", "id"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("name"),
			},
			{
				Name:        "akas",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionAkas,
				Transform:   transform.FromField("id").Transform(idToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAzureResourceGraph(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_resource_graph.listAzureResourceGraph", "session_error", err)
		return nil, err
	}

	kqlQuery := d.EqualsQuals["query"].GetStringValue()
	if kqlQuery == "" {
		return nil, nil
	}

	client := resourcegraph.NewWithBaseURI(session.ResourceManagerEndpoint)
	client.Authorizer = session.Authorizer
	ApplyRetryRules(ctx, &client, d.Connection)

	top := int32(1000)
	if d.QueryContext.Limit != nil {
		if limit := int32(*d.QueryContext.Limit); limit < top {
			top = limit
		}
	}

	options := resourcegraph.QueryRequestOptions{
		ResultFormat: resourcegraph.ResultFormatTable,
		Top:          &top,
	}

	subscriptions := []string{session.SubscriptionID}

	for {
		// Wait for rate limiting before every page fetch (including the first).
		d.WaitForListRateLimit(ctx)

		resp, err := client.Resources(ctx, resourcegraph.QueryRequest{
			Subscriptions: &subscriptions,
			Query:         &kqlQuery,
			Options:       &options,
		})
		if err != nil {
			plugin.Logger(ctx).Error("azure_resource_graph.listAzureResourceGraph", "api_error", err)
			return nil, err
		}

		rawBytes, err := json.Marshal(resp.Data)
		if err != nil {
			plugin.Logger(ctx).Error("azure_resource_graph.listAzureResourceGraph", "marshal_error", err)
			return nil, err
		}
		var table resourcegraph.Table
		if err := json.Unmarshal(rawBytes, &table); err != nil {
			plugin.Logger(ctx).Error("azure_resource_graph.listAzureResourceGraph", "unmarshal_table_error", err)
			return nil, err
		}
		if table.Columns == nil || table.Rows == nil {
			plugin.Logger(ctx).Error("azure_resource_graph.listAzureResourceGraph", "unmarshal_table_error", "response data missing columns or rows")
			return nil, fmt.Errorf("response data missing columns or rows")
		}

		// Build a slice of column names from the API-returned column descriptors.
		colNames := make([]string, len(*table.Columns))
		for i, col := range *table.Columns {
			if col.Name != nil {
				colNames[i] = *col.Name
			}
		}

		// Stream one Steampipe row per result-set row.
		for _, row := range *table.Rows {
			rowMap := make(map[string]interface{}, len(row))
			for i, cell := range row {
				if i < len(colNames) {
					rowMap[colNames[i]] = cell
				}
			}
			d.StreamListItem(ctx, rowMap)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.SkipToken == nil {
			break
		}
		options.SkipToken = resp.SkipToken
	}

	return nil, nil
}
