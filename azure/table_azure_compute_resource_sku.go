package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-09-01/skus"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureResourceSku(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_resource_sku",
		Description: "Azure Compute Resource SKU",
		List: &plugin.ListConfig{
			Hydrate: listResourceSkus,
		},

		Columns: []*plugin.Column{
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
				Transform:   transform.FromField("Sku.Capabilities"),
			},
			{
				Name:        "costs",
				Description: "A list of metadata for retrieving price info",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Sku.Costs"),
			},
			{
				Name:        "location_info",
				Description: "A list of locations and availability zones in those locations where the SKU is available",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Sku.LocationInfo"),
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
				Transform:   transform.FromField("Sku.Restrictions"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(skuDataToAkas),
			},
		},
	}
}

// custom sku struct

type skuInfo struct {
	SubscriptionID string
	Sku            skus.ResourceSku
}

//// LIST FUNCTION

func listResourceSkus(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	locksClient := skus.NewResourceSkusClient(subscriptionID)
	locksClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := locksClient.List(ctx)
		if err != nil {
			return nil, err
		}

		for _, sku := range result.Values() {
			d.StreamListItem(ctx, &skuInfo{subscriptionID, sku})
		}

		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
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
