package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureDataFactoryPipeline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_factory_pipeline",
		Description: "Azure Data Factory Pipeline",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group", "factory_name"}),
			Hydrate:           getPipeline,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate:       listPipelines,
			ParentHydrate: listFactories,
		},
		Columns: []*plugin.Column{
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
				Name:        "description",
				Description: "The description of the pipeline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Pipeline.Description"),
			},
			{
				Name:        "etag",
				Description: "Etag identifies change in the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Etag"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "factory_name",
				Description: "Time the factory was created in ISO8601 format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pipeline_folder",
				Description: "The folder that this Pipeline is in. If not specified, Pipeline will appear at the root level.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Pipeline.Folder.PipelineFolder"),
			},
			{
				Name:        "activities",
				Description: "List of activities in pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Activities"),
			},
			{
				Name:        "parameters",
				Description: "List of parameters for pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Parameters"),
			},
			{
				Name:        "variables",
				Description: "List of variables for pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Variables"),
			},
			{
				Name:        "concurrency",
				Description: "The max number of concurrent runs for the pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Concurrency"),
			},
			{
				Name:        "annotations",
				Description: "List of tags that can be used for describing the Pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Annotations"),
			},
			{
				Name:        "run_dimensions",
				Description: "Dimensions emitted by Pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.RunDimensions"),
			},
			{
				Name:        "pipeline_policy",
				Description: "Pipeline ElapsedTime Metric Policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Folder.PipelinePolicy"),
			},
			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
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

type pipelineInfo = struct {
	datafactory.PipelineResource
	FactoryName string
}

//// LIST FUNCTIONS ////

func listPipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	factoryInfo := h.Item.(datafactory.Factory)
	resourceGroup := strings.Split(*factoryInfo.ID, "/")[4]

	subscriptionID := session.SubscriptionID

	pipelineClient := datafactory.NewPipelinesClient(subscriptionID)
	pipelineClient.Authorizer = session.Authorizer
	pagesLeft := true

	for pagesLeft {
		result, err := pipelineClient.ListByFactory(ctx, resourceGroup, *factoryInfo.Name)
		if err != nil {
			return nil, err
		}

		for _, pipeline := range result.Values() {
			// plugin.Logger(ctx).Trace("getPipeline1111", pipeline)
			d.StreamListItem(ctx, pipelineInfo{pipeline, *factoryInfo.Name})
			// d.StreamListItem(ctx, pipeline)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getPipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPipeline")

	pipelineName := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	factoryName := d.KeyColumnQuals["factory_name"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	pipelineClient := datafactory.NewPipelinesClient(subscriptionID)
	pipelineClient.Authorizer = session.Authorizer

	op, err := pipelineClient.Get(ctx, resourceGroup, factoryName, pipelineName, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return pipelineInfo{op, factoryName}, nil
	}

	return op, nil
}
