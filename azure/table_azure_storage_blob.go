package azure

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	armstorage "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type blobInfo struct {
	Blob           *container.BlobItem
	Name           string
	Account        string
	Container      *string
	ResourceGroup  string
	SubscriptionID *string
	Location       string
	IsSnapshot     bool
}

//// TABLE DEFINITION

func tableAzureStorageBlob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_blob",
		Description: "Azure Storage Blob",

		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"resource_group", "storage_account_name"}),
			Hydrate:    listStorageBlobs,
			Tags: map[string]string{
				"service": "Microsoft.Storage",
				"action":  "storageAccounts/blobServices/containers/blobs/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			// Basic info
			{
				Name:        "name",
				Description: "The friendly name that identifies the blob.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_account_name",
				Description: "The friendly name that identifies the storage account, in which the blob is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "container_name",
				Description: "The friendly name that identifies the container, in which the blob has been uploaded.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Container"),
			},
			{
				Name:        "type",
				Description: "Specifies the type of the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.BlobType").Transform(derefToString),
			},
			{
				Name:        "is_snapshot",
				Description: "Specifies whether the resource is snapshot of a blob, or not.",
				Type:        proto.ColumnType_BOOL,
			},

			// Other details
			{
				Name:        "access_tier",
				Description: "The tier of the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.AccessTier").Transform(derefToString),
			},
			{
				Name:        "creation_time",
				Description: "Indicates the time, when the blob was uploaded.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.CreationTime"),
			},
			{
				Name:        "deleted",
				Description: "Specifies whether the blob was deleted, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Deleted"),
				Default:     false,
			},
			{
				Name:        "deleted_time",
				Description: "Specifies the deletion time of blob container.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.DeletedTime"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.Etag").Transform(derefToString),
			},
			{
				Name:        "last_modified",
				Description: "Specifies the date and time the container was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.LastModified"),
			},
			{
				Name:        "snapshot",
				Description: "Specifies the time, when the snapshot is taken.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Snapshot").Transform(derefToString),
			},
			{
				Name:        "version_id",
				Description: "Specifies the version id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.VersionID"),
			},
			{
				Name:        "server_encrypted",
				Description: "Indicates whether the blob is encrypted on the server, or not.",
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
				Name:        "encryption_key_sha256",
				Description: "The SHA-256 hash of the provided encryption key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CustomerProvidedKeySha256"),
			},
			{
				Name:        "is_current_version",
				Description: "Specifies whether the blob container was deleted, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.IsCurrentVersion"),
			},
			{
				Name:        "access_tier_change_time",
				Description: "Species the time, when the access tier has been updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.AccessTierChangeTime"),
			},
			{
				Name:        "access_tier_inferred",
				Description: "Indicates whether the access tier was inferred by the service.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Properties.AccessTierInferred"),
			},
			{
				Name:        "blob_sequence_number",
				Description: "Specifies the sequence number for page blob used for coordinating concurrent writes.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Blob.Properties.BlobSequenceNumber"),
			},
			{
				Name:        "content_length",
				Description: "Specifies the size of the content returned.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Blob.Properties.ContentLength"),
			},
			{
				Name:        "cache_control",
				Description: "Indicates the cache control specified for the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CacheControl"),
			},
			{
				Name:        "content_disposition",
				Description: "Specifies additional information about how to process the response payload, and also can be used to attach additional metadata.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentDisposition"),
			},
			{
				Name:        "content_encoding",
				Description: "Indicates content encoding specified for the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentEncoding"),
			},
			{
				Name:        "content_language",
				Description: "Indicates content language specified for the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentLanguage"),
			},
			{
				Name:        "content_md5",
				Description: "If the content_md5 has been set for the blob, this response header is stored so that the client can check for message content integrity.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.Properties.ContentMD5"),
			},
			{
				Name:        "content_type",
				Description: "Specifies the content type specified for the blob. If no content type was specified, the default content type is application/octet-stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ContentType"),
			},
			{
				Name:        "copy_completion_time",
				Description: "Conclusion time of the last attempted Copy Blob operation where this blob was the destination blob.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Blob.Properties.CopyCompletionTime"),
			},
			{
				Name:        "copy_id",
				Description: "A String identifier for the last attempted Copy Blob operation where this blob was the destination blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyID"),
			},
			{
				Name:        "copy_progress",
				Description: "Contains the number of bytes copied and the total bytes in the source in the last attempted Copy Blob operation where this blob was the destination blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyProgress"),
			},
			{
				Name:        "copy_source",
				Description: "An URL up to 2 KB in length that specifies the source blob used in the last attempted Copy Blob operation where this blob was the destination blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopySource"),
			},
			{
				Name:        "copy_status",
				Description: "Specifies the state of the copy operation identified by Copy ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyStatus"),
			},
			{
				Name:        "copy_status_description",
				Description: "Describes cause of fatal or non-fatal copy operation failure.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.CopyStatusDescription"),
			},
			{
				Name:        "destination_snapshot",
				Description: "Included if the blob is incremental copy blob or incremental copy snapshot, if x-ms-copy-status is success. Snapshot time of the last successful incremental copy snapshot for this blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.DestinationSnapshot"),
			},
			{
				Name:        "lease_duration",
				Description: "Specifies whether the lease is of infinite or fixed duration, when a blob is leased.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.LeaseDuration").Transform(derefToString),
			},
			{
				Name:        "lease_state",
				Description: "Specifies lease state of the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.LeaseState").Transform(derefToString),
			},
			{
				Name:        "lease_status",
				Description: "Specifies the lease status of the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.LeaseStatus").Transform(derefToString),
			},
			{
				Name:        "incremental_copy",
				Description: "Copies the snapshot of the source page blob to a destination page blob. The snapshot is copied such that only the differential changes between the previously copied snapshot are transferred to the destination.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Properties.IncrementalCopy"),
			},
			{
				Name:        "is_sealed",
				Description: "Indicate if the append blob is sealed or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Blob.Properties.IsSealed"),
			},
			{
				Name:        "remaining_retention_days",
				Description: "The number of days that the blob will be retained before being permanently deleted by the service.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Blob.Properties.RemainingRetentionDays"),
			},
			{
				Name:        "archive_status",
				Description: "Specifies the archive status of the blob.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Blob.Properties.ArchiveStatus").Transform(derefToString),
			},
			{
				Name:        "blob_tag_set",
				Description: "A list of blob tags.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.BlobTags.BlobTagSet"),
			},
			{
				Name:        "metadata",
				Description: "A name-value pair to associate with the container as metadata.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Blob.Metadata"),
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
				Transform:   transform.From(blobDataToAka),
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
				Transform:   transform.FromField("ResourceGroup").Transform(toLower),
			},
		}),
	}
}

//// FETCH FUNCTIONS

func listStorageBlobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountName := d.EqualsQuals["storage_account_name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	if accountName == "" || resourceGroup == "" {
		return nil, nil
	}

	// Management plane session (legacy) for account metadata and potential key listing
	mgmtSession, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	armCredSession, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	subscriptionID := mgmtSession.SubscriptionID

	// Get storage account properties (track2 ARM)
	armClient, err := armstorage.NewAccountsClient(subscriptionID, armCredSession.Cred, armCredSession.ClientOptions)
	if err != nil {
		return nil, err
	}
	accountResp, err := armClient.GetProperties(ctx, resourceGroup, accountName, nil)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound") || strings.Contains(err.Error(), "ResourceGroupNotFound") {
			return nil, nil
		}
		return nil, err
	}
	region := ""
	if accountResp.Account.Properties != nil && accountResp.Account.Location != nil {
		region = *accountResp.Account.Location
	}

	allowShared := true
	if accountResp.Account.Properties != nil && accountResp.Account.Properties.AllowSharedKeyAccess != nil {
		allowShared = *accountResp.Account.Properties.AllowSharedKeyAccess
	}

	config := GetConfig(d.Connection)
	authMode := "auto"
	if config.DataPlaneAuthMode != nil && *config.DataPlaneAuthMode != "" {
		authMode = strings.ToLower(*config.DataPlaneAuthMode)
	}

	// Resolve data-plane credential/client options
	serviceClient, err := buildBlobServiceClient(ctx, d, armCredSession.Cred, authMode, accountName, resourceGroup, subscriptionID, mgmtSession.StorageEndpointSuffix, allowShared)
	if err != nil {
		return nil, err
	}

	// auto fallback handled centrally inside buildBlobServiceClient

	// List containers using data-plane (track2) client
	cPager := serviceClient.NewListContainersPager(nil)
	for cPager.More() {
		// Rate limit coordination
		d.WaitForListRateLimit(ctx)
		cPage, err := cPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, c := range cPage.ContainerItems {
			containerName := ""
			if c.Name != nil {
				containerName = *c.Name
			}
			if containerName == "" {
				continue
			}
			if err := streamBlobsInContainer(ctx, d, serviceClient, accountName, resourceGroup, subscriptionID, region, containerName); err != nil {
				return nil, err
			}
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

// buildBlobServiceClient creates a track2 azblob ServiceClient according to auth_mode

func streamBlobsInContainer(ctx context.Context, d *plugin.QueryData, blobClient *azblob.Client, accountName, resourceGroup, subscriptionID, region, containerName string) error {
	pager := blobClient.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{Include: azblob.ListBlobsInclude{Copy: true, Metadata: true, Snapshots: true, Deleted: true, Tags: true, Versions: false}})
	for pager.More() {
		if d.RowsRemaining(ctx) == 0 {
			return nil
		}
		resp, err := pager.NextPage(ctx)
		if err != nil {
			// Friendly error translation
			if translated := translateBlobError(err); translated != nil {
				return translated
			}
			return err
		}
		for _, item := range resp.Segment.BlobItems {
			isSnapshot := item.Snapshot != nil && *item.Snapshot != ""
			name := ""
			if item.Name != nil {
				name = *item.Name
			}
			ci := &blobInfo{Blob: item, Name: name, Account: accountName, Container: &containerName, ResourceGroup: resourceGroup, SubscriptionID: &subscriptionID, Location: region, IsSnapshot: isSnapshot}
			d.StreamListItem(ctx, ci)
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
	}
	return nil
}

// translateBlobError defined centrally in storage_data_plane.go

//// TRANSFORM FUNCTIONS

func blobDataToAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	blob := d.HydrateItem.(*blobInfo)

	// Build resource aka
	akas := []string{"azure:///subscriptions/" + *blob.SubscriptionID + "/resourceGroups/" + blob.ResourceGroup + "/providers/Microsoft.Storage/storageAccounts/" + blob.Account + "/blobServices/default/containers/" + *blob.Container + "/blobs/" + blob.Name, "azure:///subscriptions/" + *blob.SubscriptionID + "/resourcegroups/" + strings.ToLower(blob.ResourceGroup) + "/providers/microsoft.storage/storageaccounts/" + strings.ToLower(blob.Account) + "/blobservices/default/containers/" + strings.ToLower(*blob.Container) + "/blobs/" + strings.ToLower(blob.Name)}

	return akas, nil
}

// derefToString converts pointer, fmt.Stringer, or basic values to string, returning "" when nil.
func derefToString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	v := d.Value
	if v == nil {
		return "", nil
	}
	// Unwrap pointers recursively (max a few levels to avoid cycles)
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "", nil
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return "", nil
	}
	if s, ok := rv.Interface().(fmt.Stringer); ok {
		if s == nil {
			return "", nil
		}
		return s.String(), nil
	}
	switch rv.Kind() {
	case reflect.String:
		return rv.String(), nil
	case reflect.Bool:
		if rv.Bool() {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", rv.Float()), nil
	default:
		// Fallback to fmt.Sprintf
		return fmt.Sprintf("%v", rv.Interface()), nil
	}
}
