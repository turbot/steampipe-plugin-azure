package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type queueInfo = struct {
	Queue         storage.ListQueue
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION

func tableAzureStorageQueue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_queue",
		Description: "Azure Storage Queue",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"storage_account_name", "name", "resource_group"}),
			Hydrate:    getStorageQueue,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "QueueNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageQueues,
			KeyColumns:    plugin.OptionalColumns([]string{"resource_group"}),
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the queue.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a queue uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Queue.ID"),
			},
			{
				Name:        "storage_account_name",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Queue.Type"),
			},
			{
				Name:        "metadata",
				Description: "A name-value pair that represents queue metadata.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Queue.ListQueueProperties.Metadata"),
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
				Transform:   transform.FromField("Queue.ID").Transform(idToAkas),
			},

			// Azure standard columns
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
		}),
	}
}

//// LIST FUNCTION

func listStorageQueues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	if resourceGroup != "" && resourceGroup != *account.ResourceGroup {
		return nil, nil
	}
	// Queue is not supported for the account if storage type is FileStorage or BlockBlobStorage
	if account.Account.Kind == "FileStorage" || account.Account.Kind == "BlockBlobStorage" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewQueueClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &storageClient, d.Connection)

	result, err := storageClient.List(ctx, *account.ResourceGroup, *account.Name, "", "")
	if err != nil {
		/*
		* For storage account type 'Page Blob' we are getting the kind value as 'StorageV2'.
		* Storage account type 'Page Blob' does not support table, so we are getting 'FeatureNotSupportedForAccount'/'OperationNotAllowedOnKind' error.
		* With same kind(StorageV2) of storage account, we my have different type(File Share) of storage account so we need to handle this particular error.
		 */
		if strings.Contains(err.Error(), "FeatureNotSupportedForAccount") || strings.Contains(err.Error(), "OperationNotAllowedOnKind") {
			return nil, nil
		}
		return nil, err
	}

	for _, queue := range result.Values() {
		d.StreamListItem(ctx, &queueInfo{queue, account.Name, queue.Name, account.ResourceGroup, account.Account.Location})
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
		for _, queue := range result.Values() {
			d.StreamListItem(ctx, queue)
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

func getStorageQueue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageQueue")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	accountName := d.EqualsQuals["storage_account_name"].GetStringValue()

	storageClient := storage.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &storageClient, d.Connection)

	storageDetails, err := storageClient.GetProperties(ctx, resourceGroup, accountName, "")

	if err != nil {
		return nil, err
	}

	location := storageDetails.Location

	queueClient := storage.NewQueueClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	queueClient.Authorizer = session.Authorizer

	op, err := queueClient.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	if op.QueueProperties != nil {
		return &queueInfo{
			Queue: storage.ListQueue{
				Name: op.Name,
				ID:   op.ID,
				Type: op.Type,
				ListQueueProperties: &storage.ListQueueProperties{
					Metadata: op.QueueProperties.Metadata,
				},
			},
			Account:       &accountName,
			Name:          &name,
			ResourceGroup: &resourceGroup,
			Location:      location,
		}, nil
	}

	return nil, nil
}
