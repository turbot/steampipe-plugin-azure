package azure

import (
	"context"
	"strings"

	sub "github.com/Azure/azure-sdk-for-go/profiles/latest/subscription/mgmt/subscription"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeVirtualMachineSize(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_virtual_machine_size",
		Description: "Azure Compute Virtual Machine Size",
		List: &plugin.ListConfig{
			ParentHydrate: listLocations,
			Hydrate:       listComputeVirtualMachineSizes,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the virtual machine size.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineSize.Name"),
			},
			{
				Name:        "number_of_cores",
				Description: "The number of cores supported by the virtual machine size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VirtualMachineSize.NumberOfCores"),
			},
			{
				Name:        "os_disk_size_in_mb",
				Description: "The OS disk size, in MB, allowed by the virtual machine size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VirtualMachineSize.OSDiskSizeInMB"),
			},
			{
				Name:        "resource_disk_size_in_mb",
				Description: "The resource disk size, in MB, allowed by the virtual machine size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VirtualMachineSize.ResourceDiskSizeInMB"),
			},
			{
				Name:        "memory_in_mb",
				Description: "The amount of memory, in MB, supported by the virtual machine size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VirtualMachineSize.MemoryInMB"),
			},
			{
				Name:        "max_data_disk_count",
				Description: "The maximum number of data disks that can be attached to the virtual machine size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VirtualMachineSize.MaxDataDiskCount"),
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineSize.Name"),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
		}),
	}
}

type VMSizeInfo struct {
	*armcompute.VirtualMachineSize
	Location string
}

//// LIST FUNCTION ////

func listComputeVirtualMachineSizes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	if h.Item == nil {
		return nil, nil
	}

	locationDetails := h.Item.(sub.Location)

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_compute_virtual_machine_size.listComputeVirtualMachineSizes", "session_error", err)
		return nil, err
	}

	session.ClientOptions.APIVersion = "2024-07-01"

	clientFactory, err := armcompute.NewVirtualMachineSizesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	pager := clientFactory.NewListPager(*locationDetails.Name, &armcompute.VirtualMachineSizesClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			// In Azure, resource providers are services that allow you to interact with resources (like virtual machines). 
			// The relevant resource provider for Virtual Machines is `Microsoft.Compute`. 
			// If this provider is not registered, or if it's not available in the specified region, you might encounter the error.
                        // You can use the command (`az provider show --namespace Microsoft.Compute`) to check the availability of the service in the specified location.
                        // Look for the `locations/vmSizes` resource type in the command result to verify it's availability.
			
			if strings.Contains(strings.ToLower(err.Error()), "no registered resource provider found for location") {
				plugin.Logger(ctx).Error("azure_compute_virtual_machine_size.listComputeVirtualMachineSizes", "no registered resource provider found", err.Error())
				return nil, nil
			}
			plugin.Logger(ctx).Error("azure_compute_virtual_machine_size.listComputeVirtualMachineSizes", "api_error", err)
			return nil, err
		}

		for _, virtualMachineSize := range page.Value {
			d.StreamListItem(ctx, VMSizeInfo{virtualMachineSize, *locationDetails.Name})
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

