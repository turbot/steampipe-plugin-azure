package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeAvailabilitySet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_availability_set",
		Description: "Azure Compute Availability Set",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAzureComputeAvailabilitySet,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureComputeAvailabilitySets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the availability set",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource in Azure",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_fault_domain_count",
				Description: "Contains the fault domain count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AvailabilitySetProperties.PlatformFaultDomainCount"),
			},
			{
				Name:        "platform_update_domain_count",
				Description: "Contains the update domain count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AvailabilitySetProperties.PlatformUpdateDomainCount"),
			},
			{
				Name:        "proximity_placement_group_id",
				Description: "Specifies information about the proximity placement group that the availability set should be assigned to",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AvailabilitySetProperties.ProximityPlacementGroup.ID"),
			},
			{
				Name:        "sku_capacity",
				Description: "Specifies the number of virtual machines in the scale set",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_name",
				Description: "The availability sets sku name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "Specifies the tier of virtual machines in a scale set",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "status",
				Description: "The resource status information",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AvailabilitySetProperties.Statuses"),
			},
			{
				Name:        "virtual_machines",
				Description: "A list of references to all virtual machines in the availability set",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AvailabilitySetProperties.VirtualMachines"),
			},

			// Standard columns
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION ////

func listAzureComputeAvailabilitySets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeAvailabilitySets")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewAvailabilitySetsClient(subscriptionID)
	client.Authorizer = session.Authorizer
	pagesLeft := true

	for pagesLeft {
		result, err := client.ListBySubscription(context.Background(), "")
		if err != nil {
			return nil, err
		}

		for _, availabilitySet := range result.Values() {
			d.StreamListItem(ctx, availabilitySet)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS ////

func getAzureComputeAvailabilitySet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeAvailabilitySet")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewAvailabilitySetsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(context.Background(), resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
