package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureComputeImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_image",
		Description: "Azure Compute Image",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getComputeImage,
			Tags: map[string]string{
				"service": "compute",
				"action":  "Microsoft.Compute/images/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeImages,
			Tags: map[string]string{
				"service": "compute",
				"action":  "Microsoft.Compute/images/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the image",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a image uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the image resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.ProvisioningState"),
			},
			{
				Name:        "hyper_v_generation",
				Description: "Gets the HyperVGenerationType of the VirtualMachine created from the image",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.HyperVGeneration").Transform(transform.ToString),
			},
			{
				Name:        "source_virtual_machine_id",
				Description: "Contains the id of the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.SourceVirtualMachine.ID"),
			},
			{
				Name:        "storage_profile_os_disk_blob_uri",
				Description: "Contains uri of the virtual hard disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.BlobURI"),
			},
			{
				Name:        "storage_profile_os_disk_caching",
				Description: "Specifies the caching requirements",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.Caching").Transform(transform.ToString),
			},
			{
				Name:        "storage_profile_os_disk_encryption_set",
				Description: "Specifies the customer managed disk encryption set resource id for the managed image disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.DiskEncryptionSet.ID"),
			},
			{
				Name:        "storage_profile_os_disk_managed_disk_id",
				Description: "Contains the id of the managed disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.ManagedDisk.ID"),
			},
			{
				Name:        "storage_profile_os_disk_size_gb",
				Description: "Specifies the size of empty data disks in gigabytes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.DiskSizeGB"),
			},
			{
				Name:        "storage_profile_os_disk_snapshot_id",
				Description: "Contains the id of the snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.Snapshot.ID"),
			},
			{
				Name:        "storage_profile_os_disk_storage_account_type",
				Description: "Specifies the storage account type for the managed disk",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.StorageAccountType").Transform(transform.ToString),
			},
			{
				Name:        "storage_profile_os_disk_state",
				Description: "Contains state of the OS",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.OsState").Transform(transform.ToString),
			},
			{
				Name:        "storage_profile_os_disk_type",
				Description: "Specifies the type of the OS that is included in the disk if creating a VM from a custom image",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageProperties.StorageProfile.OsDisk.OsType").Transform(transform.ToString),
			},
			{
				Name:        "storage_profile_zone_resilient",
				Description: "Specifies whether an image is zone resilient or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ImageProperties.StorageProfile.ZoneResilient"),
				Default:     false,
			},
			{
				Name:        "storage_profile_data_disks",
				Description: "A list of parameters that are used to add a data disk to a virtual machine",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ImageProperties.StorageProfile.DataDisks"),
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
				Transform:   transform.FromField("Tags"),
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

//// FETCH FUNCTIONS ////

func listComputeImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	computeClient := compute.NewImagesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	computeClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &computeClient, d.Connection)

	// Apply rate limiting
	d.WaitForListRateLimit(ctx)

	result, err := computeClient.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, image := range result.Values() {
		d.StreamListItem(ctx, image)
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

		for _, image := range result.Values() {
			d.StreamListItem(ctx, image)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getComputeImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeImage")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	computeClient := compute.NewImagesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	computeClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &computeClient, d.Connection)

	op, err := computeClient.Get(ctx, resourceGroup, name, "")
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
