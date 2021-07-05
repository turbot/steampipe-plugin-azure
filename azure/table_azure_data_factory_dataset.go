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

func tableAzureDataFactoryDataset(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_factory_dataset",
		Description: "Azure Data Factory Dataset",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group", "factory_name"}),
			Hydrate:           getDataset,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate:       listDatasets,
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
				Description: "The description of the Dataset.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Dataset.Description"),
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
				Name:        "properties",
				Description: "Dataset ElapsedTime Metric Policy.",
				Type:        proto.ColumnType_JSON,
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

type DatasetInfo = struct {
	datafactory.DatasetResource
	FactoryName string
}

//// LIST FUNCTIONS ////

func listDatasets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	factoryInfo := h.Item.(datafactory.Factory)
	resourceGroup := strings.Split(*factoryInfo.ID, "/")[4]

	subscriptionID := session.SubscriptionID

	datasetClient := datafactory.NewDatasetsClient(subscriptionID)
	datasetClient.Authorizer = session.Authorizer
	pagesLeft := true

	for pagesLeft {
		result, err := datasetClient.ListByFactory(ctx, resourceGroup, *factoryInfo.Name)
		if err != nil {
			return nil, err
		}

		for _, dataset := range result.Values() {
			d.StreamListItem(ctx, DatasetInfo{dataset, *factoryInfo.Name})
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getDataset(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataset")

	DatasetName := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	factoryName := d.KeyColumnQuals["factory_name"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	datasetClient := datafactory.NewDatasetsClient(subscriptionID)
	datasetClient.Authorizer = session.Authorizer

	op, err := datasetClient.Get(ctx, resourceGroup, factoryName, DatasetName, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return DatasetInfo{op, factoryName}, nil
	}

	return nil, nil
}
