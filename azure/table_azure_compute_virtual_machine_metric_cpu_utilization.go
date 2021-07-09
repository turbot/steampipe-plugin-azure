package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureComputeInstanceCpuUtilizationMetric(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_virtual_machine_metric_cpu_utilization",
		Description: "Azure Compute Virtual Machine Metrics - CPU Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listAzureComputeVirtualMachines,
			Hydrate:       listComputeVirtualMachineMetricCpuUtilization,
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

func listComputeVirtualMachineMetricCpuUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vmInfo := h.Item.(compute.VirtualMachine)

	return listAzureMonitorMetricStatistics(ctx, d, "FIVE_MINUTES", "Microsoft.Compute/virtualMachines", "Percentage CPU", *vmInfo.ID)
}
