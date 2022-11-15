package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_snapshot",
		Description: "Azure Compute Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAzureComputeSnapshot,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureComputeSnapshots,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the snapshot",
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
			},
			{
				Name:        "create_option",
				Description: "Specifies the possible sources of a disk's creation",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.CreationData.CreateOption").Transform(transform.ToString),
			},
			{
				Name:        "disk_access_id",
				Description: "ARM id of the DiskAccess resource for using private endpoints on disks",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.DiskAccessID"),
			},
			{
				Name:        "disk_encryption_set_id",
				Description: "ResourceId of the disk encryption set to use for enabling encryption at rest",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.Encryption.DiskEncryptionSetID"),
			},
			{
				Name:        "disk_size_bytes",
				Description: "The size of the disk in bytes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SnapshotProperties.DiskSizeBytes"),
			},
			{
				Name:        "disk_size_gb",
				Description: "The size of the disk to create",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SnapshotProperties.DiskSizeGB"),
			},
			{
				Name:        "encryption_setting_collection_enabled",
				Description: "Specifies whether the encryption is enables, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SnapshotProperties.EncryptionSettingsCollection.Enabled"),
			},
			{
				Name:        "encryption_setting_version",
				Description: "Describes what type of encryption is used for the disks",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.EncryptionSettingsCollection.EncryptionSettingsVersion"),
			},
			{
				Name:        "encryption_type",
				Description: "The type of the encryption",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.Encryption.Type").Transform(transform.ToString),
			},
			{
				Name:        "gallery_image_reference_id",
				Description: "A relative uri containing either a Platform Image Repository or user image reference",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.CreationData.GalleryImageReference.ID"),
			},
			{
				Name:        "gallery_reference_lun",
				Description: "Specifies the index that indicates which of the data disks in the image to use",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SnapshotProperties.CreationData.GalleryImageReference.Lun"),
			},
			{
				Name:        "hyperv_generation",
				Description: "Specifies the hypervisor generation of the Virtual Machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.HyperVGeneration").Transform(transform.ToString),
			},
			{
				Name:        "image_reference_id",
				Description: "A relative uri containing either a Platform Image Repository or user image reference",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.CreationData.ImageReference.ID"),
			},
			{
				Name:        "image_reference_lun",
				Description: "Specifies the index that indicates which of the data disks in the image to use",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SnapshotProperties.CreationData.ImageReference.Lun"),
			},
			{
				Name:        "incremental",
				Description: "Specifies whether a snapshot is incremental, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SnapshotProperties.Incremental"),
			},
			{
				Name:        "network_access_policy",
				Description: "Contains the type of access",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.NetworkAccessPolicy").Transform(transform.ToString),
			},
			{
				Name:        "os_type",
				Description: "Contains the type of operating system",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.OsType").Transform(transform.ToString),
			},
			{
				Name:        "sku_name",
				Description: "The snapshot sku name",
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
				Name:        "source_resource_id",
				Description: "ARM id of the source snapshot or disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.CreationData.SourceResourceID"),
			},
			{
				Name:        "source_unique_id",
				Description: "An unique id identifying the source of this resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.CreationData.SourceUniqueID"),
			},
			{
				Name:        "source_uri",
				Description: "An URI of a blob to be imported into a managed disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.CreationData.SourceURI"),
			},
			{
				Name:        "storage_account_id",
				Description: "The Azure Resource Manager identifier of the storage account containing the blob to import as a disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.CreationData.StorageAccountID"),
			},
			{
				Name:        "time_created",
				Description: "The time when the snapshot was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SnapshotProperties.TimeCreated").Transform(convertDateToTime),
			},
			{
				Name:        "unique_id",
				Description: "An unique Guid identifying the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotProperties.UniqueID"),
			},
			{
				Name:        "upload_size_bytes",
				Description: "The size of the contents of the upload including the VHD footer",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SnapshotProperties.CreationData.UploadSizeBytes"),
			},
			{
				Name:        "encryption_settings",
				Description: "A list of encryption settings, one for each disk volume",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SnapshotProperties.EncryptionSettingsCollection.EncryptionSettings"),
			},
			{
				Name:        "virtual_machines",
				Description: "A list of references to all virtual machines in the availability set",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AvailabilitySetProperties.VirtualMachines"),
			},

			// Steampipe standard columns
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION ////

func listAzureComputeSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeSnapshots")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewSnapshotsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, snapshot := range result.Values() {
		d.StreamListItem(ctx, snapshot)
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

		for _, snapshot := range result.Values() {
			d.StreamListItem(ctx, snapshot)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS ////

func getAzureComputeSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeSnapshot")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewSnapshotsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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
