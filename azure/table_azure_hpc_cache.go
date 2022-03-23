package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-05-01/storagecache"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureHPCCache(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_hpc_cache",
		Description: "Azure HPC Cache",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getHPCCache,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listHPCCaches,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID of the cache.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "ARM provisioning state. Possible values include: 'Succeeded', 'Failed', 'Cancelled', 'Creating', 'Deleting', 'Updating'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CacheProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the cache.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_size_gb",
				Description: "The size of the cache, in GB.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CacheProperties.CacheSizeGB"),
			},
			{
				Name:        "sku_name",
				Description: "The SKU for the cache.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "subnet",
				Description: "Subnet used for the cache.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CacheProperties.Subnet"),
			},
			{
				Name:        "directory_services_settings",
				Description: "Specifies directory services settings of the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CacheProperties.DirectoryServicesSettings"),
			},
			{
				Name:        "encryption_settings",
				Description: "Specifies encryption settings of the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CacheProperties.EncryptionSettings"),
			},
			{
				Name:        "health",
				Description: "The health of the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CacheProperties.Health"),
			},
			{
				Name:        "identity",
				Description: "The identity of the cache, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "mount_addresses",
				Description: "Array of IP addresses that can be used by clients mounting the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CacheProperties.MountAddresses"),
			},
			{
				Name:        "network_settings",
				Description: "Specifies network settings of the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractHPCCacheNetworkSettings),
			},
			{
				Name:        "security_settings",
				Description: "Specifies security settings of the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CacheProperties.SecuritySettings"),
			},
			{
				Name:        "system_data",
				Description: "The system meta data relating to the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "upgrade_status",
				Description: "Upgrade status of the cache.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractHPCCacheUpgradeStatus),
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

type CacheNetworkSettingsInfo struct {
	Mtu              *int32
	UtilityAddresses *[]string
	DNSServers       *[]string
	DNSSearchDomain  *string
	NtpServer        *string
}

type CacheUpgradeStatusInfo struct {
	CurrentFirmwareVersion *string
	FirmwareUpdateStatus   interface{}
	FirmwareUpdateDeadline *date.Time
	LastFirmwareUpdate     *date.Time
	PendingFirmwareVersion *string
}

//// LIST FUNCTION

func listHPCCaches(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := storagecache.NewCachesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listHPCCaches", "list", err)
		return nil, err
	}

	for _, cache := range result.Values() {
		d.StreamListItem(ctx, cache)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listHPCCaches", "list_paging", err)
			return nil, err
		}
		for _, cache := range result.Values() {
			d.StreamListItem(ctx, cache)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getHPCCache(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHPCCache")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := storagecache.NewCachesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getHPCCache", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

// If we return the API response directly, the output does not provide
// all the properties of NetworkSettings
func extractHPCCacheNetworkSettings(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	cache := d.HydrateItem.(storagecache.Cache)
	var properties CacheNetworkSettingsInfo

	if cache.CacheProperties.NetworkSettings != nil {
		if cache.CacheProperties.NetworkSettings.Mtu != nil {
			properties.Mtu = cache.CacheProperties.NetworkSettings.Mtu
		}
		if cache.CacheProperties.NetworkSettings.UtilityAddresses != nil {
			properties.UtilityAddresses = cache.CacheProperties.NetworkSettings.UtilityAddresses
		}
		if cache.CacheProperties.NetworkSettings.DNSServers != nil {
			properties.DNSServers = cache.CacheProperties.NetworkSettings.DNSServers
		}
		if cache.CacheProperties.NetworkSettings.DNSSearchDomain != nil {
			properties.DNSSearchDomain = cache.CacheProperties.NetworkSettings.DNSSearchDomain
		}
		if cache.CacheProperties.NetworkSettings.NtpServer != nil {
			properties.NtpServer = cache.CacheProperties.NetworkSettings.NtpServer
		}
	}

	return properties, nil
}

// If we return the API response directly, the output does not provide
// all the properties of UpgradeStatus
func extractHPCCacheUpgradeStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	cache := d.HydrateItem.(storagecache.Cache)
	var properties CacheUpgradeStatusInfo

	if cache.CacheProperties.UpgradeStatus != nil {
		if cache.CacheProperties.UpgradeStatus.CurrentFirmwareVersion != nil {
			properties.CurrentFirmwareVersion = cache.CacheProperties.UpgradeStatus.CurrentFirmwareVersion
		}
		if len(cache.CacheProperties.UpgradeStatus.FirmwareUpdateStatus) > 0 {
			properties.FirmwareUpdateStatus = cache.CacheProperties.UpgradeStatus.FirmwareUpdateStatus
		}
		if cache.CacheProperties.UpgradeStatus.FirmwareUpdateDeadline != nil {
			properties.FirmwareUpdateDeadline = cache.CacheProperties.UpgradeStatus.FirmwareUpdateDeadline
		}
		if cache.CacheProperties.UpgradeStatus.LastFirmwareUpdate != nil {
			properties.LastFirmwareUpdate = cache.CacheProperties.UpgradeStatus.LastFirmwareUpdate
		}
		if cache.CacheProperties.UpgradeStatus.PendingFirmwareVersion != nil {
			properties.PendingFirmwareVersion = cache.CacheProperties.UpgradeStatus.PendingFirmwareVersion
		}
	}

	return properties, nil
}
