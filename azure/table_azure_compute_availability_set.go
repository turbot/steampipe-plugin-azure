package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeAvailabilitySet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_availability_set",
		Description: "Azure Compute Availability Set",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAzureComputeAvailabilitySet,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureComputeAvailabilitySets,
		},
		Columns: azureColumns([]*plugin.Column{
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
				Hydrate:     getAzureComputeAvailabilitySet,
				Transform:   transform.From(extractStatusForAvailabilitySet),
			},
			{
				Name:        "virtual_machines",
				Description: "A list of references to all virtual machines in the availability set",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAzureComputeAvailabilitySet,
				Transform:   transform.From(extractVirtualMachinesForAvailabilitySet),
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

//// LIST FUNCTION ////

func listAzureComputeAvailabilitySets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeAvailabilitySets")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewAvailabilitySetsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	result, err := client.ListBySubscription(ctx, "")
	if err != nil {
		return nil, err
	}

	for _, availabilitySet := range result.Values() {
		d.StreamListItem(ctx, availabilitySet)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err := result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, availabilitySet := range result.Values() {
			d.StreamListItem(ctx, availabilitySet)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS ////

func getAzureComputeAvailabilitySet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	resourceGroup, name := "", ""

	name = d.EqualsQuals["name"].GetStringValue()
	resourceGroup = d.EqualsQuals["resource_group"].GetStringValue()

	if h.Item != nil {
		availabilitySet := h.Item.(compute.AvailabilitySet)
		id := availabilitySet.ID
		resourceGroup = strings.Split(*id, "/")[4]
		name = *availabilitySet.Name
	}

	// Empty check
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewAvailabilitySetsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.Get(ctx, resourceGroup, name)
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

//// UTILITY FUNCTION

func extractVirtualMachinesForAvailabilitySet(_ context.Context, d *transform.TransformData) (interface{}, error) {
	availabilitySet := d.HydrateItem.(compute.AvailabilitySet)
	var properties []map[string]interface{}

	if availabilitySet.AvailabilitySetProperties != nil && availabilitySet.AvailabilitySetProperties.VirtualMachines != nil {
		vmProperies := availabilitySet.AvailabilitySetProperties
		for _, i := range *vmProperies.VirtualMachines {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			properties = append(properties, objectMap)
		}
	}

	return properties, nil
}

func extractStatusForAvailabilitySet(_ context.Context, d *transform.TransformData) (interface{}, error) {
	availabilitySet := d.HydrateItem.(compute.AvailabilitySet)
	var properties []map[string]interface{}

	if availabilitySet.AvailabilitySetProperties != nil && availabilitySet.AvailabilitySetProperties.Statuses != nil {
		properies := availabilitySet.AvailabilitySetProperties
		for _, i := range *properies.Statuses {
			objectMap := make(map[string]interface{})
			if i.Code != nil {
				objectMap["code"] = i.Code
			}
			if i.DisplayStatus != nil {
				objectMap["displayStatus"] = i.DisplayStatus
			}
			if i.Level != "" {
				objectMap["level"] = i.Level
			}
			if i.Message != nil {
				objectMap["message"] = i.Message
			}

			properties = append(properties, objectMap)
		}
	}

	return properties, nil
}
