package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureComputeDisksReadOpsMetric(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_disk_metric_read_ops",
		Description: "Azure Compute Disk Metrics - Read Ops",
		List: &plugin.ListConfig{
			ParentHydrate: listAzureComputeDisks,
			Hydrate:       listComputeDiskMetricReadOps,
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

func listComputeDiskMetricReadOps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	diskInfo := h.Item.(compute.Disk)

	return listAzureMonitorMetricStatistics(ctx, d, "FIVE_MINUTES", "Microsoft.Compute/disks", "Composite Disk Read Operations/sec", *diskInfo.ID)
}
