package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type blobServiceInfo = struct {
	Blob          storage.BlobServiceProperties
	Account       *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION ////

func tableAzureStorageBlobService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_blob_service",
		Description: "Azure Storage Blob Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"storage_account_name", "resource_group"}),
			Hydrate:    getStorageBlobService,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageBlobServices,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the blob",
				Transform:   transform.FromField("Blob.Name"),
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a blob uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.ID"),
			},
			{
				Name:        "storage_account_name",
				Description: "A unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Type"),
			},
			{
				Name:        "sku_name",
				Description: "The sku name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Contains the sku tier",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "automatic_snapshot_policy_enabled",
				Description: "Specifies whether automatic snapshot creation is enabled, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.AutomaticSnapshotPolicyEnabled"),
				Default:     false,
			},
			{
				Name:        "change_feed_enabled",
				Description: "Specifies whether change feed event logging is enabled for the Blob service",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.ChangeFeed.Enabled"),
				Default:     false,
			},
			{
				Name:        "default_service_version",
				Description: "Indicates the default version to use for requests to the Blob service if an incoming request’s version is not specified",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.DefaultServiceVersion"),
			},
			{
				Name:        "is_versioning_enabled",
				Description: "Specifies whether the versioning is enabled, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.IsVersioningEnabled"),
				Default:     false,
			},
			{
				Name:        "container_delete_retention_policy",
				Description: "The blob service properties for container soft delete",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.ContainerDeleteRetentionPolicy"),
			},
			{
				Name:        "cors_rules",
				Description: "A list of CORS rules for a storage account’s Blob service",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.Cors.CorsRules"),
			},
			{
				Name:        "delete_retention_policy",
				Description: "The blob service properties for blob soft delete",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.DeleteRetentionPolicy"),
			},
			{
				Name:        "restore_policy",
				Description: "The blob service properties for blob restore policy",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.BlobServicePropertiesProperties.RestorePolicy"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.ID").Transform(idToAkas),
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
				Transform:   transform.FromField("Blob.ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listStorageBlobServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	// Blob is not supported for the account if storage type is FileStorage
	if account.Account.Kind == "FileStorage" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewBlobServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	result, err := storageClient.List(ctx, *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	for _, blobService := range *result.Value {
		d.StreamListItem(ctx, &blobServiceInfo{blobService, account.Name, account.ResourceGroup, account.Account.Location})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getStorageBlobService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageBlobService")

	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	accountName := d.EqualsQuals["storage_account_name"].GetStringValue()

	// length of the AccountName must be greater than, or equal to 3, and
	// length of the ResourceGroupName must be greater than 1
	if len(accountName) < 3 || len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	storageDetails, err := storageClient.GetProperties(ctx, resourceGroup, accountName, "")

	if err != nil {
		return nil, err
	}

	location := storageDetails.Location

	blobClient := storage.NewBlobServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	blobClient.Authorizer = session.Authorizer

	op, err := blobClient.GetServiceProperties(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}

	return &blobServiceInfo{op, &accountName, &resourceGroup, location}, nil
}
