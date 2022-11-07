package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-09-01/skus"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAzureResourceSku(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_resource_sku",
		Description: "Azure Compute Resource SKU",
		List: &plugin.ListConfig{
			Hydrate: listResourceSkus,
		},

		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "resource_type",
				Description: "The type of resource the SKU applies to",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.ResourceType"),
			},
			{
				Name:        "tier",
				Description: "Specifies the tier of virtual machines in a scale set",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "size",
				Description: "The Size of the SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Size"),
			},
			{
				Name:        "family",
				Description: "The Family of this particular SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Family"),
			},
			{
				Name:        "kind",
				Description: "The Kind of resources that are supported in this SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Kind"),
			},
			{
				Name:        "default_capacity",
				Description: "Contains the default capacity",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity.Default"),
			},
			{
				Name:        "maximum_capacity",
				Description: "The maximum capacity that can be set",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity.Maximum"),
			},
			{
				Name:        "minimum_capacity",
				Description: "The minimum capacity that can be set",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity.Minimum"),
			},
			{
				Name:        "scale_type",
				Description: "The scale type applicable to the sku",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Capacity.ScaleType").Transform(transform.ToString),
			},
			{
				Name:        "api_versions",
				Description: "The api versions that support this SKU",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Sku.APIVersions"),
			},
			{
				Name:        "capabilities",
				Description: "A name value pair to describe the capability",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromMethod("ComputeResourceSkuCapabilities"),
			},
			{
				Name:        "costs",
				Description: "A list of metadata for retrieving price info",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromMethod("ComputeResourceSkuCosts"),
			},
			{
				Name:        "location_info",
				Description: "A list of locations and availability zones in those locations where the SKU is available",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromMethod("ComputeResourceSkuLocationInfo"),
			},
			{
				Name:        "locations",
				Description: "The set of locations that the SKU is available",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Sku.Locations"),
			},
			{
				Name:        "restrictions",
				Description: "The restrictions because of which SKU cannot be used",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromMethod("ComputeResourceSkuRestrictions"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(skuDataToAkas),
			},
		}),
	}
}

// custom sku struct

type skuInfo struct {
	SubscriptionID string
	Sku skus.ResourceSku
}

//// LIST FUNCTION

func listResourceSkus(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	locksClient := skus.NewResourceSkusClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	locksClient.Authorizer = session.Authorizer

	result, err := locksClient.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, sku := range result.Values() {
		d.StreamListItem(ctx, &skuInfo{subscriptionID, sku})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, sku := range result.Values() {
			d.StreamListItem(ctx, &skuInfo{subscriptionID, sku})
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// TRANSFORM FUNCTION ////

func skuDataToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	sku := d.HydrateItem.(*skuInfo)
	locations := *sku.Sku.Locations
	id := "azure:///subscriptions/" + sku.SubscriptionID + "/locations/" + locations[0] + "/resourcetypes" + *sku.Sku.ResourceType + "name/" + *sku.Sku.Name
	akas := []string{strings.ToLower(id)}
	return akas, nil
}

//// HELPER TRANSFORM FUNCTIONS to populate columns always returning [{}]

func (skuData *skuInfo) ComputeResourceSkuCapabilities() []map[string]interface{} {
	if skuData.Sku.Capabilities == nil {
		return nil
	}
	capabilities := []map[string]interface{}{}

	for _, a := range *skuData.Sku.Capabilities {
		data := map[string]interface{}{}
		if a.Name != nil {
			data["name"] = *a.Name
		}
		if a.Value != nil {
			data["value"] = *a.Value
		}
		capabilities = append(capabilities, data)
	}

	return capabilities
}

func (skuData *skuInfo) ComputeResourceSkuRestrictions() []map[string]interface{} {
	if skuData.Sku.Capabilities == nil {
		return nil
	}
	restrictions := []map[string]interface{}{}

	for _, a := range *skuData.Sku.Restrictions {
		data := map[string]interface{}{}
		data["type"] = &a.Type
		data["reasonCode"] = &a.ReasonCode
		if a.Values != nil {
			data["Values"] = *a.Values
		}
		if a.RestrictionInfo != nil {
			data["restrictionInfo"] = *a.RestrictionInfo
		}
		restrictions = append(restrictions, data)
	}

	return restrictions
}

func (skuData *skuInfo) ComputeResourceSkuLocationInfo() []map[string]interface{} {
	if skuData.Sku.LocationInfo == nil {
		return nil
	}
	locationInfo := []map[string]interface{}{}

	for _, a := range *skuData.Sku.LocationInfo {
		data := map[string]interface{}{}
		if a.Location != nil {
			data["location"] = *a.Location
		}
		if a.Zones != nil {
			data["zones"] = *a.Zones
		}
		locationInfo = append(locationInfo, data)
	}

	return locationInfo
}

func (skuData *skuInfo) ComputeResourceSkuCosts() []map[string]interface{} {
	if skuData.Sku.Costs == nil {
		return nil
	}
	costs := []map[string]interface{}{}

	for _, a := range *skuData.Sku.Costs {
		data := map[string]interface{}{}
		if a.MeterID != nil {
			data["meterID"] = *a.MeterID
		}
		if a.Quantity != nil {
			data["quantity"] = *a.Quantity
		}
		if a.ExtendedUnit != nil {
			data["extendedUnit"] = *a.ExtendedUnit
		}
		costs = append(costs, data)
	}

	return costs
}
