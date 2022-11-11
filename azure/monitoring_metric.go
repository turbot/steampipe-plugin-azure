package azure

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2022-10-01-preview/insights"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type monitoringMetric struct {
	// Resource Name
	DimensionValue string
	// MetadataValue represents a metric metadata value.
	MetaData *insights.MetadataValue
	// Metric the result data of a query.
	Metric *insights.Metric
	// The maximum metric value for the data point.
	Maximum *float64
	// The minimum metric value for the data point.
	Minimum *float64
	// The average of the metric values that correspond to the data point.
	Average *float64
	// The number of metric values that contributed to the aggregate value of this data point.
	SampleCount *float64
	// The sum of the metric values for the data point.
	Sum *float64
	// The time stamp used for the data point.
	TimeStamp string
	// The units in which the metric value is reported.
	Unit string
}

//// TABLE DEFINITION

func monitoringMetricColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonMonitoringMetricColumns()...)
}

func commonMonitoringMetricColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "maximum",
			Description: "The maximum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "minimum",
			Description: "The minimum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "average",
			Description: "The average of the metric values that correspond to the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sample_count",
			Description: "The number of metric values that contributed to the aggregate value of this data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sum",
			Description: "The sum of the metric values for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "timestamp",
			Description: "The time stamp used for the data point.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("TimeStamp"),
		},
		{
			Name:        "unit",
			Description: "The units in which the metric value is reported.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "cloud_environment",
			Description: ColumnDescriptionCloudEnvironment,
			Type:        proto.ColumnType_STRING,
			Hydrate:     plugin.HydrateFunc(getCloudEnvironment).WithCache(),
			Transform:   transform.FromValue(),
		},
		{
			Name:        "resource_group",
			Description: ColumnDescriptionResourceGroup,
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("DimensionValue").Transform(extractResourceGroupFromID),
		},
		{
			Name:        "subscription_id",
			Description: ColumnDescriptionSubscription,
			Type:        proto.ColumnType_STRING,
			Hydrate:     plugin.HydrateFunc(getSubscriptionID).WithCache(),
			Transform:   transform.FromValue(),
		},
	}
}

func getMonitoringIntervalForGranularity(granularity string) string {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 24 hours
		return "PT24H"
	case "HOURLY":
		// 1 hour
		return "PT1H"
	}
	// else 5 minutes
	return "PT5M"
}

func getMonitoringStartDateForGranularity(granularity string) string {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// Last 1 year
		return time.Now().UTC().AddDate(-1, 0, 0).Format(time.RFC3339)
	case "HOURLY":
		// Last 60 days
		return time.Now().UTC().AddDate(0, 0, -60).Format(time.RFC3339)
	}
	// Last 5 days
	return time.Now().UTC().AddDate(0, 0, -5).Format(time.RFC3339)
}

func listAzureMonitorMetricStatistics(ctx context.Context, d *plugin.QueryData, granularity string, metricNameSpace string, metricNames string, dimensionValue string) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	monitoringClient := insights.NewMetricsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	monitoringClient.Authorizer = session.Authorizer

	// Define param values
	interval := getMonitoringIntervalForGranularity(granularity)
	aggregation := "average,count,maximum,minimum,total"
	timeSpan := getMonitoringStartDateForGranularity(granularity) + "/" + time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339) // Retrieve data within a year
	orderBy := "timestamp"
	top := int32(1000) // Maximum number of record fetch with given interval
	filter := ""

	result, err := monitoringClient.List(ctx, dimensionValue, timeSpan, &interval, metricNames, aggregation, &top, orderBy, filter, insights.ResultTypeData, metricNameSpace)
	if err != nil {
		return nil, err
	}
	for _, metric := range *result.Value {
		for _, timeseries := range *metric.Timeseries {
			for _, data := range *timeseries.Data {
				if data.Average != nil {
					d.StreamListItem(ctx, &monitoringMetric{
						DimensionValue: dimensionValue,
						TimeStamp:      data.TimeStamp.Format(time.RFC3339),
						Maximum:        data.Maximum,
						Minimum:        data.Minimum,
						Average:        data.Average,
						Sum:            data.Total,
						SampleCount:    data.Count,
						Unit:           string(metric.Unit),
					})
				}
			}
		}
	}

	return nil, nil
}
