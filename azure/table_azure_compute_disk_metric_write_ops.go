package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableComputeDisksWriteOpsMetric(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_disk_metric_write_ops",
		Description: "Azure Compute Disk Metrics - Write Ops",
		List: &plugin.ListConfig{
			ParentHydrate: listAzureComputeDisks,
			Hydrate:       listComputeDiskMetricWriteOps,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the disk.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue").Transform(lastPathElement),
			},
		}),
	}
}

//// LIST FUNCTION

func listComputeDiskMetricWriteOps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	diskInfo := h.Item.(compute.Disk)

	return listAzureMonitorMetricStatistics(ctx, d, "FIVE_MINUTES", "Microsoft.Compute/disks", "Composite Disk Write Operations/sec", *diskInfo.ID)
}
