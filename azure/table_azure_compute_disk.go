package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeDisk(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_disk",
		Description: "Azure Compute Disk",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAzureComputeDisk,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureComputeDisks,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the disk",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource in Azure",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The disk provisioning state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.ProvisioningState"),
			},
			{
				Name:        "managed_by",
				Description: "A relative URI containing the ID of the VM that has the disk attached",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku_name",
				Description: "The disks sku name. Can be Standard_LRS, Premium_LRS, StandardSSD_LRS, or UltraSSD_LRS",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "The sku tier",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "time_created",
				Description: "The time when the disk was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DiskProperties.TimeCreated").Transform(convertDateToTime),
			},
			{
				Name:        "unique_id",
				Description: "Unique Guid identifying the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.UniqueID"),
			},
			{
				Name:        "disk_access_id",
				Description: "ARM id of the DiskAccess resource for using private endpoints on disks",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.DiskAccessID"),
			},
			{
				Name:        "disk_size_bytes",
				Description: "The size of the disk in bytes",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("DiskProperties.DiskSizeBytes"),
			},
			{
				Name:        "disk_size_gb",
				Description: "If this field is present for updates or creation with other options, it indicates a resize. Resizes are only allowed if the disk is not attached to a running VM, and can only increase the disk's size",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DiskProperties.DiskSizeGB"),
			},
			{
				Name:        "disk_state",
				Description: "This enumerates the possible state of the disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.DiskState"),
			},
			{
				Name:        "hyper_v_generation",
				Description: "The hypervisor generation of the Virtual Machine. Applicable to OS disks only",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.HyperVGeneration"),
			},
			{
				Name:        "disk_iops_read_only",
				Description: "The total number of IOPS that will be allowed across all VMs mounting the shared disk as ReadOnly. One operation can transfer between 4k and 256k bytes",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("DiskProperties.DiskIOPSReadOnly"),
			},
			{
				Name:        "disk_iops_read_write",
				Description: "The number of IOPS allowed for this disk; only settable for UltraSSD disks. One operation can transfer between 4k and 256k bytes",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("DiskProperties.DiskIOPSReadWrite"),
			},
			{
				Name:        "disk_iops_mbps_read_only",
				Description: "The total throughput (MBps) that will be allowed across all VMs mounting the shared disk as ReadOnly. MBps means millions of bytes per second - MB here uses the ISO notation, of powers of 10",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DiskProperties.DiskMBpsReadOnly"),
			},
			{
				Name:        "disk_iops_mbps_read_write",
				Description: "The bandwidth allowed for this disk; only settable for UltraSSD disks. MBps means millions of bytes per second - MB here uses the ISO notation, of powers of 10",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DiskProperties.DiskMBpsReadWrite"),
			},
			{
				Name:        "max_shares",
				Description: "The maximum number of VMs that can attach to the disk at the same time. Value greater than one indicates a disk that can be mounted on multiple VMs at the same time",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DiskProperties.MaxShares"),
			},
			{
				Name:        "os_type",
				Description: "The Operating System type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.OsType"),
			},
			{
				Name:        "encryption_settings_collection_enabled",
				Description: "Shows the status of the encryption settings for the disk",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("DiskProperties.EncryptionSettingsCollection.Enabled"),
			},
			{
				Name:        "encryption_settings_collection_version",
				Description: "Describes the type of encryption is used for the disks. '1.0' corresponds to Azure Disk Encryption with AAD app. '1.1' corresponds to Azure Disk Encryption",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.EncryptionSettingsCollection.EncryptionSettingsVersion"),
			},
			{
				Name:        "encryption_disk_encryption_set_id",
				Description: "ResourceId of the disk encryption set to use for enabling encryption at rest",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.Encryption.DiskEncryptionSetID"),
			},
			{
				Name:        "encryption_type",
				Description: "The type of key used to encrypt the data of the disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.Encryption.Type"),
			},
			{
				Name:        "network_access_policy",
				Description: "Policy for accessing the disk via network",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_data_option",
				Description: "This enumerates the possible sources of a disk's creation",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.CreationData.CreateOption"),
			},
			{
				Name:        "creation_data_storage_account_id",
				Description: "The Azure Resource Manager identifier of the storage account containing the blob to import as a disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.CreationData.StorageAccountID"),
			},
			{
				Name:        "creation_data_source_uri",
				Description: "The URI of a blob to be imported into a managed disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.CreationData.SourceURI"),
			},
			{
				Name:        "creation_data_source_resource_id",
				Description: "The ARM id of the source snapshot or disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.CreationData.SourceResourceID"),
			},
			{
				Name:        "creation_data_source_unique_id",
				Description: "An unique id identifying the source of this resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.CreationData.SourceUniqueID"),
			},
			{
				Name:        "creation_data_upload_size_bytes",
				Description: "This is the size of the contents of the upload including the VHD footer. This value should be between 20972032 (20 MiB + 512 bytes for the VHD footer) and 35183298347520 bytes (32 TiB + 512 bytes for the VHD footer)",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DiskProperties.CreationData.UploadSizeBytes"),
			},
			{
				Name:        "creation_data_image_reference_id",
				Description: "A relative uri containing either a Platform Image Repository or user image reference",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.CreationData.ImageReference.ID"),
			},
			{
				Name:        "creation_data_image_reference_lun",
				Description: "If the disk is created from an image's data disk, this is an index that indicates which of the data disks in the image to use. For OS disks, this field is null",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DiskProperties.CreationData.ImageReference.Lun"),
			},
			{
				Name:        "creation_data_gallery_image_reference_id",
				Description: "The ARM id of the shared galley image version from which disk was created",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskProperties.CreationData.GalleryImageReference.ID"),
			},
			{
				Name:        "creation_data_gallery_image_reference_lun",
				Description: "An index that indicates which of the data disks in the image to use, if the disk is created from an image's data disk",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DiskProperties.CreationData.GalleryImageReference.Lun"),
			},
			{
				Name:        "encryption_settings_collection_settings",
				Description: "A collection of encryption settings, one for each disk volume",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DiskProperties.EncryptionSettingsCollection.EncryptionSettings"),
			},
			{
				Name:        "share_info",
				Description: "Details of the list of all VMs that have the disk attached",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DiskProperties.ShareInfo"),
			},
			{
				Name:        "zones",
				Description: "The Logical zone list for Disk",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "managed_by_extended",
				Description: "List of relative URIs containing the IDs of the VMs that have the disk attached",
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
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
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

//// LIST FUNCTION ////

func listAzureComputeDisks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeDisks")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewDisksClient(subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, disk := range result.Values() {
		d.StreamListItem(ctx, disk)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, disk := range result.Values() {
			d.StreamListItem(ctx, disk)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION ////

func getAzureComputeDisk(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeDisk")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewDisksClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
