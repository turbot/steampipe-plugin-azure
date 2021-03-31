package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
)

//// TABLE DEFINITION

func tableAzureStorageContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_container",
		Description: "Azure Storage Container",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group", "account_name"}),
			Hydrate:           getStorageContainer,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageContainers,
		},

		Columns: []*plugin.Column{
			// Basic info
			{
				Name:        "name",
				Description: "The friendly name that identifies the container.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The container ID",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "type",
				Description: "Specifies the type of the container.",
				Type:        proto.ColumnType_STRING,
			},
			// Other details
			{
				Name:        "account_name",
				Description: "The friendly name that identifies the storage account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToAccountName),
			},
			{
				Name:        "container_properties",
				Description: "The blob container properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerProperties"),
			},

			// Standard steampipe columns
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
				Transform:   transform.FromField("Etag"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Standard azure columns
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

//// FETCH FUNCTIONS

func listStorageContainers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	// List all containers
	containerClient := storage.NewBlobContainersClient(subscriptionID)
	containerClient.Authorizer = session.Authorizer
	pagesLeft := true
	for pagesLeft {
		containerList, err := containerClient.List(ctx, *account.ResourceGroup, *account.Name, "", "", "")
		if err != nil {
			return nil, err
		}

		for _, container := range containerList.Values() {
			d.StreamLeafListItem(ctx, container)
		}
		containerList.NextWithContext(context.Background())
		pagesLeft = containerList.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getStorageContainer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageContainer")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	accountName := d.KeyColumnQuals["account_name"].GetStringValue()

	storageClient := storage.NewBlobContainersClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer

	op, err := storageClient.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
