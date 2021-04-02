package azure

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
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
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "ContainerNotFound"}),
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
				Transform:   transform.FromGo(),
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
				Name:        "deleted",
				Description: "Indicates whether the blob container was deleted.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ContainerProperties.Deleted"),
			},
			{
				Name:        "default_encryption_scope",
				Description: "Default the container to use specified encryption scope for all writes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.DefaultEncryptionScope"),
			},
			{
				Name:        "public_access",
				Description: "Specifies whether data in the container may be accessed publicly and the level of access.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.PublicAccess"),
			},
			{
				Name:        "remaining_retention_days",
				Description: "Remaining retention days for soft deleted blob container.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ContainerProperties.RemainingRetentionDays"),
			},
			{
				Name:        "version",
				Description: "The version of the deleted blob container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.Version"),
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

//// LIST FUNCTION

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

//// TRANSFORM FUNCTIONS

func idToAccountName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	accountName := strings.Split(id, "/")[8]
	return accountName, nil
}
