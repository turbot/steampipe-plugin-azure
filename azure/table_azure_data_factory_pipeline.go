package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/datafactory/mgmt/datafactory"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureDataFactoryPipeline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_factory_pipeline",
		Description: "Azure Data Factory Pipeline",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group", "factory_name"}),
			Hydrate:    getDataFactoryPipeline,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate:       listDataFactoryPipelines,
			ParentHydrate: listDataFactories,
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
				Name:        "factory_name",
				Description: "Name of the factory the pipeline belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the pipeline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Pipeline.Description"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "concurrency",
				Description: "The max number of concurrent runs for the pipeline.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Pipeline.Concurrency"),
			},
			{
				Name:        "pipeline_folder",
				Description: "The folder that this Pipeline is in. If not specified, Pipeline will appear at the root level.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Pipeline.Folder.PipelineFolder"),
			},
			{
				Name:        "activities",
				Description: "A list of activities in pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Activities"),
			},
			{
				Name:        "annotations",
				Description: "A list of tags that can be used for describing the Pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Annotations"),
			},
			{
				Name:        "parameters",
				Description: "A list of parameters for pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Parameters"),
			},
			{
				Name:        "pipeline_policy",
				Description: "Pipeline ElapsedTime Metric Policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Folder.PipelinePolicy"),
			},
			{
				Name:        "variables",
				Description: "A list of variables for pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Variables"),
			},
			{
				Name:        "run_dimensions",
				Description: "Dimensions emitted by Pipeline.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.RunDimensions"),
			},

			// Steampipe standard columns
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

			// Azure standard columns
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

type pipelineInfo = struct {
	datafactory.PipelineResource
	FactoryName string
}

//// LIST FUNCTION

func listDataFactoryPipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	// Get factory details
	factoryInfo := h.Item.(datafactory.Factory)
	resourceGroup := strings.Split(*factoryInfo.ID, "/")[4]

	pipelineClient := datafactory.NewPipelinesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	pipelineClient.Authorizer = session.Authorizer

	result, err := pipelineClient.ListByFactory(ctx, resourceGroup, *factoryInfo.Name)
	if err != nil {
		return nil, err
	}
	for _, pipeline := range result.Values() {
		d.StreamListItem(ctx, pipelineInfo{pipeline, *factoryInfo.Name})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, pipeline := range result.Values() {
			d.StreamListItem(ctx, pipelineInfo{pipeline, *factoryInfo.Name})
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDataFactoryPipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataFactoryPipeline")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	pipelineClient := datafactory.NewPipelinesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	pipelineClient.Authorizer = session.Authorizer

	pipelineName := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	factoryName := d.EqualsQuals["factory_name"].GetStringValue()

	// Return nil, if no input provided
	if pipelineName == "" || resourceGroup == "" || factoryName == "" {
		return nil, nil
	}

	op, err := pipelineClient.Get(ctx, resourceGroup, factoryName, pipelineName, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return pipelineInfo{op, factoryName}, nil
	}

	return nil, nil
}
