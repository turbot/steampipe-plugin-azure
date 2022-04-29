package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureComputeDiskMetricWriteOpsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_disk_metric_write_ops_daily",
		Description: "Azure Compute Disk Metrics - Write Ops (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listAzureComputeDisks,
			Hydrate:       listComputeDiskMetricWriteOpsDaily,
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

func listComputeDiskMetricWriteOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	diskInfo := h.Item.(compute.Disk)

	return listAzureMonitorMetricStatistics(ctx, d, "DAILY", "Microsoft.Compute/disks", "Composite Disk Write Operations/sec", *diskInfo.ID)
}
