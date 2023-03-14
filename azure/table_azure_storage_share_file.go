package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureStorageShareFile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_share_file",
		Description: "Azure Storage Share File",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group", "storage_account_name"}),
			Hydrate:    getStorageAccountsFileShare,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageAccountsFileShares,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource.",
			},
			{
				Name:        "storage_account_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the storage account.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Fully qualified resource ID for the resource.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the resource.",
			},
			{
				Name:        "access_tier",
				Type:        proto.ColumnType_STRING,
				Description: "Access tier for specific share. GpV2 account can choose between TransactionOptimized (default), Hot, and Cool.",
				Transform:   transform.FromField("FileShareProperties.AccessTier"),
			},
			{
				Name:        "access_tier_change_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Indicates the last modification time for share access tier.",
				Transform:   transform.FromField("FileShareProperties.AccessTierChangeTime").Transform(convertDateToTime),
			},
			{
				Name:        "access_tier_status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates if there is a pending transition for access tier.",
				Transform:   transform.FromField("FileShareProperties.AccessTierStatus"),
			},
			{
				Name:        "last_modified_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Returns the date and time the share was last modified.",
				Transform:   transform.FromField("FileShareProperties.LastModifiedTime").Transform(convertDateToTime),
			},
			{
				Name:        "deleted",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the share was deleted.",
				Transform:   transform.FromField("FileShareProperties.Deleted"),
			},
			{
				Name:        "deleted_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The deleted time if the share was deleted.",
				Transform:   transform.FromField("FileShareProperties.DeletedTime").Transform(convertDateToTime),
			},
			{
				Name:        "enabled_protocols",
				Type:        proto.ColumnType_STRING,
				Description: "The authentication protocol that is used for the file share. Can only be specified when creating a share. Possible values include: 'SMB', 'NFS'.",
				Transform:   transform.FromField("FileShareProperties.EnabledProtocols"),
			},
			{
				Name:        "remaining_retention_days",
				Type:        proto.ColumnType_INT,
				Description: "Remaining retention days for share that was soft deleted.",
				Transform:   transform.FromField("FileShareProperties.RemainingRetentionDays"),
			},
			{
				Name:        "root_squash",
				Type:        proto.ColumnType_STRING,
				Description: "The property is for NFS share only. The default is NoRootSquash. Possible values include: 'NoRootSquash', 'RootSquash', 'AllSquash'.",
				Transform:   transform.FromField("FileShareProperties.RootSquash"),
			},
			{
				Name:        "share_quota",
				Type:        proto.ColumnType_INT,
				Description: "The maximum size of the share, in gigabytes. Must be greater than 0, and less than or equal to 5TB (5120). For Large File Shares, the maximum size is 102400.",
				Transform:   transform.FromField("FileShareProperties.ShareQuota"),
			},
			{
				Name:        "share_usage_bytes",
				Type:        proto.ColumnType_INT,
				Description: "The approximate size of the data stored on the share. Note that this value may not include all recently created or recently resized files.",
				Transform:   transform.FromField("FileShareProperties.ShareUsageBytes"),
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "The version of the share.",
				Transform:   transform.FromField("FileShareProperties.Version"),
			},
			{
				Name:        "metadata",
				Type:        proto.ColumnType_JSON,
				Description: "A name-value pair to associate with the share as metadata.",
				Transform:   transform.FromField("FileShareProperties.Metadata"),
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
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceGroup").Transform(toLower),
			},
		}),
	}
}

type FileShareInfo struct {
	storage.FileShareProperties
	Name               string
	ID                 string
	Type               string
	StorageAccountName string
	ResourceGroup      string
}

//// LIST FUNCTION

func listStorageAccountsFileShares(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	storageAccount := h.Item.(*storageAccountInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	logger := plugin.Logger(ctx)
	if err != nil {
		logger.Error("listStorageAccountsFileShare", "get session error", err)
		return nil, err
	}

	if storageAccount.Account.Kind == "BlockBlobStorage" {
		return nil, nil
	}

	subscriptionID := session.SubscriptionID
	fileShareCLient := storage.NewFileSharesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	fileShareCLient.Authorizer = session.Authorizer

	// Limiting the results
	limit := d.QueryContext.Limit
	maxResult := "100"
	if d.QueryContext.Limit != nil {
		if *limit < 100 {
			maxResult = types.IntToString(*limit)
		}
	}

	result, err := fileShareCLient.List(ctx, *storageAccount.ResourceGroup, *storageAccount.Name, maxResult, "", "")
	if err != nil {
		logger.Error("listStorageAccountsFileShare", "api error", err)

		// This api throws FeatureNotSupportedForAccount or OperationNotAllowedOnKind error if the storage account kind is not File Share
		if strings.Contains(err.Error(), "FeatureNotSupportedForAccount") || strings.Contains(err.Error(), "OperationNotAllowedOnKind") {
			return nil, nil
		}

		return nil, err
	}

	for _, fileShare := range result.Values() {
		d.StreamListItem(ctx, &FileShareInfo{
			FileShareProperties: *fileShare.FileShareProperties,
			Name:                *fileShare.Name,
			ID:                  *fileShare.ID,
			Type:                *fileShare.Type,
			StorageAccountName:  *storageAccount.Name,
			ResourceGroup:       *storageAccount.ResourceGroup,
		})

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

		for _, fileShare := range result.Values() {
			d.StreamListItem(ctx, &FileShareInfo{
				FileShareProperties: *fileShare.FileShareProperties,
				Name:                *fileShare.Name,
				ID:                  *fileShare.ID,
				Type:                *fileShare.Type,
				StorageAccountName:  *storageAccount.Name,
				ResourceGroup:       *storageAccount.ResourceGroup,
			})

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

func getStorageAccountsFileShare(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageAccountsFileShare")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	logger := plugin.Logger(ctx)
	if err != nil {
		logger.Error("getStorageAccountsFileShare", "get session error", err)
		return nil, err
	}

	resourceGroup := d.EqualsQualString("resource_group")
	storageAccountName := d.EqualsQualString("storage_account_name")
	name := d.EqualsQualString("name")

	if strings.Trim(name, " ") == "" || strings.Trim(resourceGroup, " ") == "" || strings.Trim(storageAccountName, " ") == "" {
		return nil, nil
	}

	subscriptionID := session.SubscriptionID
	fileShareCLient := storage.NewFileSharesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	fileShareCLient.Authorizer = session.Authorizer

	result, err := fileShareCLient.Get(ctx, resourceGroup, storageAccountName, name, "")

	if err != nil {
		return nil, err
	}

	return &FileShareInfo{
		FileShareProperties: *result.FileShareProperties,
		Name:                *result.Name,
		ID:                  *result.ID,
		Type:                *result.Type,
		StorageAccountName:  storageAccountName,
		ResourceGroup:       resourceGroup,
	}, nil
}
