package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureStreamAnalyticsJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_stream_analytics_job",
		Description: "Azure Stream Analytics Job",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getStreamAnalyticsJob,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listStreamAnalyticsJobs,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "job_id",
				Description: "A GUID uniquely identifying the streaming job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.JobID"),
			},
			{
				Name:        "job_state",
				Description: "Describes the state of the streaming job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.JobState"),
			},
			{
				Name:        "provisioning_state",
				Description: "Describes the provisioning status of the streaming job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compatibility_level",
				Description: "Controls certain runtime behaviors of the streaming job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.CompatibilityLevel"),
			},
			{
				Name:        "created_date",
				Description: "Specifies the time when the stream analytics job was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StreamingJobProperties.CreatedDate").Transform(convertDateToTime),
			},
			{
				Name:        "data_locale",
				Description: "The data locale of the stream analytics job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.DataLocale"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.Etag"),
			},
			{
				Name:        "events_late_arrival_max_delay_in_seconds",
				Description: "The maximum tolerable delay in seconds where events arriving late could be included.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("StreamingJobProperties.EventsLateArrivalMaxDelayInSeconds"),
			},
			{
				Name:        "events_out_of_order_max_delay_in_seconds",
				Description: "The maximum tolerable delay in seconds where out-of-order events can be adjusted to be back in order.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("StreamingJobProperties.EventsOutOfOrderMaxDelayInSeconds"),
			},
			{
				Name:        "events_out_of_order_policy",
				Description: "Indicates the policy to apply to events that arrive out of order in the input event stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.EventsOutOfOrderPolicy"),
			},
			{
				Name:        "last_output_event_time",
				Description: "Indicating the last output event time of the streaming job or null indicating that output has not yet been produced.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StreamingJobProperties.LastOutputEventTime").Transform(convertDateToTime),
			},
			{
				Name:        "output_error_policy",
				Description: "Indicates the policy to apply to events that arrive at the output and cannot be written to the external storage due to being malformed (missing column values, column values of wrong type or size).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.OutputErrorPolicy"),
			},
			{
				Name:        "output_start_mode",
				Description: "This property should only be utilized when it is desired that the job be started immediately upon creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.OutputStartMode"),
			},
			{
				Name:        "output_start_time",
				Description: "Indicates the starting point of the output event stream, or null to indicate that the output event stream will start whenever the streaming job is started.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StreamingJobProperties.OutputStartTime").Transform(convertDateToTime),
			},
			{
				Name:        "sku_name",
				Description: "Describes the sku name of the streaming job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamingJobProperties.Sku.Name"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the streaming job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listStreamAnalyticsJobDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "functions",
				Description: "A list of one or more functions for the streaming job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StreamingJobProperties.functions"),
			},
			{
				Name:        "inputs",
				Description: "A list of one or more inputs to the streaming job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StreamingJobProperties.Inputs"),
			},
			{
				Name:        "outputs",
				Description: "A list of one or more outputs for the streaming job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StreamingJobProperties.Outputs"),
			},
			{
				Name:        "transformation",
				Description: "Indicates the query and the number of streaming units to use for the streaming job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StreamingJobProperties.Transformation"),
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
				Transform:   transform.FromField("Location").Transform(formatRegion).Transform(toLower),
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

//// LIST FUNCTION

func listStreamAnalyticsJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	streamingJobsClient := streamanalytics.NewStreamingJobsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	streamingJobsClient.Authorizer = session.Authorizer

	result, err := streamingJobsClient.List(context.Background(), "")
	if err != nil {
		return nil, err
	}
	for _, streamingJob := range result.Values() {
		d.StreamListItem(ctx, streamingJob)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, streamingJob := range result.Values() {
			d.StreamListItem(ctx, streamingJob)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getStreamAnalyticsJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStrestreamingJobamAnalyticsJob")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	streamingJobsClient := streamanalytics.NewStreamingJobsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	streamingJobsClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := streamingJobsClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listStreamAnalyticsJobDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listStreamAnalyticsJobDiagnosticSettings")
	id := *h.Item.(streamanalytics.StreamingJob).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["id"] = i.ID
		}
		if i.Name != nil {
			objectMap["name"] = i.Name
		}
		if i.Type != nil {
			objectMap["type"] = i.Type
		}
		if i.DiagnosticSettings != nil {
			objectMap["properties"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}
	return diagnosticSettings, nil
}
