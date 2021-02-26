package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type queueInfo = struct {
	Queue         storage.ListQueue
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION ////

func tableAzureStorageQueue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_queue",
		Description: "Azure Storage Queue",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"storage_account_name", "name", "resource_group"}),
			ItemFromKey:       queueDataFromKey,
			Hydrate:           getStorageQueue,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "QueueNotFound", "ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageQueues,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the queue",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a queue uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Queue.ID"),
			},
			{
				Name:        "storage_account_name",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Queue.Type"),
			},
			{
				Name:        "metadata",
				Description: "A name-value pair that represents queue metadata",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Queue.ListQueueProperties.Metadata"),
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
				Transform:   transform.FromField("Queue.ID").Transform(idToAkas),
			},
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Queue.ID").Transform(extractResourceGroupFromID),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Queue.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// BUILD HYDRATE INPUT ////

func queueDataFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	resourceGroup := quals["resource_group"].GetStringValue()
	accountName := quals["storage_account_name"].GetStringValue()
	item := &queueInfo{
		Account:       &accountName,
		Name:          &name,
		ResourceGroup: &resourceGroup,
	}
	return item, nil
}

//// FETCH FUNCTIONS ////

func listStorageQueues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewQueueClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer

	result, err := storageClient.List(context.Background(), *account.ResourceGroup, *account.Name, "", "")
	if err != nil {
		return nil, err
	}

	for _, queue := range result.Values() {
		d.StreamLeafListItem(ctx, &queueInfo{queue, account.Name, queue.Name, account.ResourceGroup, account.Account.Location})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getStorageQueue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	queueData := h.Item.(*queueInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewAccountsClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer

	storageDetails, err := storageClient.GetProperties(context.Background(), *queueData.ResourceGroup, *queueData.Account, "")

	if err != nil {
		return nil, err
	}

	location := storageDetails.Location

	queueClient := storage.NewQueueClient(subscriptionID)
	queueClient.Authorizer = session.Authorizer

	op, err := queueClient.Get(context.Background(), *queueData.ResourceGroup, *queueData.Account, *queueData.Name)
	if err != nil {
		return nil, err
	}

	if op.QueueProperties != nil {
		queueData = &queueInfo{
			Queue: storage.ListQueue{
				Name: op.Name,
				ID:   op.ID,
				Type: op.Type,
				ListQueueProperties: &storage.ListQueueProperties{
					Metadata: op.QueueProperties.Metadata,
				},
			},
			Account:       queueData.Account,
			Name:          queueData.Name,
			ResourceGroup: queueData.ResourceGroup,
			Location:      location,
		}
		return queueData, nil
	}

	return nil, nil
}
