package azure

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
)

type blobInfo = struct {
	Blob           azblob.BlobItemInternal
	Name           string
	Account        *string
	Container      *string
	ResourceGroup  *string
	SubscriptionID *string
	Location       *string
	IsSnapshot     bool
}

//// TABLE DEFINITION

func tableAzureStorageBlob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_blob",
		Description: "Azure Storage Blob",
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageBlobs,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the blob",
			},
			{
				Name:        "storage_account_name",
				Description: "A unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "container_name",
				Description: "A unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Container"),
			},
			{
				Name:        "type",
				Description: "Type of the blob",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.BlobType").Transform(transform.ToString),
			},
			{
				Name:        "is_snapshot",
				Description: "Specifies whether the resource is snapshot of a blob, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "access_tier",
				Description: "The tier of the blob",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.AccessTier").Transform(transform.ToString),
			},
			{
				Name:        "etag",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.Etag").Transform(transform.ToString),
			},
			{
				Name:        "deleted",
				Description: "Specifies whether the blob container was deleted, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Deleted"),
			},
			{
				Name:        "snapshot",
				Description: "Specifies the snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Snapshot"),
			},
			{
				Name:        "destination_snapshot",
				Description: "Specifies the snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.DestinationSnapshot"),
			},
			{
				Name:        "creation_time",
				Description: "Contains ID to identify a blob uniquely",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.CreationTime"),
			},
			{
				Name:        "version_id",
				Description: "Specifies the version id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.VersionID"),
			},
			{
				Name:        "server_encrypted",
				Description: "Checks whether there are at least one existing tag, or not. If set to false, all existing legal hold tags are cleared out.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Properties.ServerEncrypted"),
			},
			{
				Name:        "encryption_scope",
				Description: "The name of the encryption scope under which the blob is encrypted.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.EncryptionScope"),
			},
			{
				Name:        "is_current_version",
				Description: "Specifies whether the blob container was deleted, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.IsCurrentVersion"),
			},
			{
				Name:        "deleted_time",
				Description: "Specifies the deletion time of blob container.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.DeletedTime"),
			},
			{
				Name:        "last_modified",
				Description: "Specifies the date and time the container was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.LastModified"),
			},
			{
				Name:        "access_tier_change_time",
				Description: "Default the container to use specified encryption scope for all writes.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.AccessTierChangeTime"),
			},
			{
				Name:        "access_tier_inferred",
				Description: "Block override of encryption scope from the container default.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Properties.AccessTierInferred"),
			},
			{
				Name:        "blob_sequence_number",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Blob.Properties.BlobSequenceNumber"),
			},
			{
				Name:        "content_length",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Blob.Properties.ContentLength"),
			},
			{
				Name:        "cache_control",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CacheControl"),
			},
			{
				Name:        "content_disposition",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentDisposition"),
			},
			{
				Name:        "content_encoding",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentEncoding"),
			},
			{
				Name:        "content_language",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentLanguage"),
			},
			{
				Name:        "content_md5",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.Properties.ContentMD5"),
			},
			{
				Name:        "content_type",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentType"),
			},
			{
				Name:        "copy_completion_time",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.CopyCompletionTime"),
			},
			{
				Name:        "copy_id",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyID"),
			},
			{
				Name:        "copy_progress",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyProgress"),
			},
			{
				Name:        "copy_source",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopySource"),
			},
			{
				Name:        "copy_status",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyStatus"),
			},
			{
				Name:        "copy_status_description",
				Description: "Indicates whether any ImmutabilityPolicy has been created for this container, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyStatusDescription"),
			},
			{
				Name:        "lease_duration",
				Description: "Specifies whether the lease on a container is of infinite or fixed duration, only when the container is leased.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.LeaseDuration").Transform(transform.ToString),
			},
			{
				Name:        "lease_state",
				Description: "Specifies lease state of the container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.LeaseState").Transform(transform.ToString),
			},
			{
				Name:        "lease_status",
				Description: "Specifies lease status of the container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.LeaseStatus").Transform(transform.ToString),
			},
			{
				Name:        "incremental_copy",
				Description: "Specifies whether data in the container may be accessed publicly and the level of access.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Properties.IncrementalCopy"),
			},
			{
				Name:        "is_sealed",
				Description: "Specifies whether data in the container may be accessed publicly and the level of access.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Properties.IsSealed"),
			},
			{
				Name:        "remaining_retention_days",
				Description: "Specifies whether data in the container may be accessed publicly and the level of access.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Blob.Properties.RemainingRetentionDays"),
			},
			{
				Name:        "archive_status",
				Description: "The version of the deleted blob container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ArchiveStatus").Transform(transform.ToString),
			},
			{
				Name:        "blob_tag_set",
				Description: "The LegalHold property of the container.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.BlobTags.BlobTagSet"),
			},
			{
				Name:        "metadata",
				Description: "A name-value pair to associate with the container as metadata.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.Metadata"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(blobDataToAka),
			},
			{
				Name:        "region",
				Description: "The Azure region in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "resource_group",
				Description: "Name of the resource group, the blob is created at",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_id",
				Description: "The Azure Subscription ID in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubscriptionID"),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listStorageBlobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	// List storage account keys
	storageClient := storage.NewAccountsClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer
	keys, err := storageClient.ListKeys(ctx, *account.ResourceGroup, *account.Name, "")
	if err != nil {
		return nil, err
	}

	credential, errC := azblob.NewSharedKeyCredential(*account.Account.Name, *(*(keys.Keys))[0].Value)
	if errC != nil {
		return nil, errC
	}

	// List all containers
	containerClient := storage.NewBlobContainersClient(subscriptionID)
	containerClient.Authorizer = session.Authorizer
	var containers []storage.ListContainerItem
	pagesLeft := true
	for pagesLeft {
		containerList, err := containerClient.List(ctx, *account.ResourceGroup, *account.Name, "", "", "")
		if err != nil {
			return nil, err
		}

		for _, container := range containerList.Values() {
			containers = append(containers, container)
		}
		containerList.NextWithContext(context.Background())
		pagesLeft = containerList.NotDone()
	}

	var wg sync.WaitGroup
	blobCh := make(chan []blobInfo, len(containers))
	errorCh := make(chan error, len(containers))

	// Iterating all the available containers
	for _, item := range containers {
		wg.Add(1)
		go getRowDataForBlobAsync(ctx, item, account, credential, &wg, blobCh, errorCh)
	}

	// wait for all containers to be processed
	wg.Wait()

	// NOTE: close channel before ranging over results
	close(blobCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		return nil, err
	}

	for item := range blobCh {
		for _, data := range item {
			d.StreamLeafListItem(ctx, &blobInfo{data.Blob, data.Name, account.Name, data.Container, account.ResourceGroup, &subscriptionID, account.Account.Location, data.IsSnapshot})
		}
	}

	return nil, err
}

func getRowDataForBlobAsync(ctx context.Context, item storage.ListContainerItem, account *storageAccountInfo, credential *azblob.SharedKeyCredential, wg *sync.WaitGroup, subnetCh chan []blobInfo, errorCh chan error) {
	defer wg.Done()

	rowData, err := getRowDataForBlob(ctx, item, account, credential)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		subnetCh <- rowData
	}
}

// List all the available blobs
func getRowDataForBlob(ctx context.Context, container storage.ListContainerItem, account *storageAccountInfo, credential *azblob.SharedKeyCredential) ([]blobInfo, error) {
	primaryURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", *account.Name))
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create Service URL
	serviceURL := azblob.NewServiceURL(*primaryURL, p)
	containerURL := serviceURL.NewContainerURL(*container.Name)

	var items []blobInfo
	subscriptionID := strings.Split(string(*container.ID), "/")[2]

	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{
			Details: azblob.BlobListingDetails{
				Copy:             false,
				Metadata:         false,
				Snapshots:        true,
				UncommittedBlobs: false,
				Deleted:          false,
				Tags:             false,
				Versions:         false,
			},
		})
		if err != nil {
			return nil, err
		}

		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker
		isSnapshot := true

		for _, blob := range listBlob.Segment.BlobItems {
			// Snapshot of a blob has same configuration,
			// only difference is that the snapshot has a property which specifies
			// the time, when the snapshot was taken
			if len(blob.Snapshot) < 1 {
				isSnapshot = false
			}
			items = append(items, blobInfo{blob, blob.Name, account.Name, container.Name, account.ResourceGroup, &subscriptionID, account.Account.Location, isSnapshot})
		}
	}

	return items, nil
}

//// TRANSFORM FUNCTIONS

func blobDataToAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	blob := d.HydrateItem.(*blobInfo)

	// /subscriptions/{subscription-id}/resourceGroups/res9871/providers/Microsoft.Storage/storageAccounts/sto6217/blobServices/default/containers/container1634
	akas := []string{"azure:///subscriptions/" + *blob.SubscriptionID + "/resourceGroups/" + *blob.ResourceGroup + "/providers/Microsoft.Storage/storageAccounts/" + *blob.Account + "/blobServices/default/containers/" + *blob.Container + "/blobs/" + blob.Name, "azure:///subscriptions/" + *blob.SubscriptionID + "/resourcegroups/" + strings.ToLower(*blob.ResourceGroup) + "/providers/microsoft.storage/storageaccounts/" + strings.ToLower(*blob.Account) + "/blobservices/default/containers/" + strings.ToLower(*blob.Container) + "/blobs/" + strings.ToLower(blob.Name)}

	return akas, nil
}
