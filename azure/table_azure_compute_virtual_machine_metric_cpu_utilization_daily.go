package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureComputeVirtualMachineMetricCpuUtilizationDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_virtual_machine_metric_cpu_utilization_daily",
		Description: "Azure Compute Virtual Machine Metrics - CPU Utilization (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeVirtualMachines,
			Hydrate:       listComputeVirtualMachineMetricCpuUtilizationDaily,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue").Transform(lastPathElement),
			},
		}),
	}
}

//// LIST FUNCTION

func listComputeVirtualMachineMetricCpuUtilizationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vmInfo := h.Item.(compute.VirtualMachine)

	return listAzureMonitorMetricStatistics(ctx, d, "DAILY", "Microsoft.Compute/virtualMachines", "Percentage CPU", *vmInfo.ID)
}
