package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/accounts"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type storageAccountInfo = struct {
	Account       storage.Account
	Name          *string
	ResourceGroup *string
}

//// TABLE DEFINITION

func tableAzureStorageAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_account",
		Description: "Azure Storage Account",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getStorageAccount,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listStorageAccounts,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the storage account.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a storage account uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.ID"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Type"),
			},
			{
				Name:        "access_tier",
				Description: "The access tier used for billing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.AccessTier").Transform(transform.ToString),
			},
			{
				Name:        "kind",
				Description: "The kind of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Kind").Transform(transform.ToString),
			},
			{
				Name:        "sku_name",
				Description: "Contains sku name of the storage account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Contains sku tier of the storage account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "creation_time",
				Description: "Creation date and time of the storage account.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Account.AccountProperties.CreationTime").Transform(convertDateToTime),
			},
			{
				Name:        "allow_blob_public_access",
				Description: "Specifies whether allow or disallow public access to all blobs or containers in the storage account.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Account.AccountProperties.AllowBlobPublicAccess"),
			},
			{
				Name:        "blob_change_feed_enabled",
				Description: "Specifies whether change feed event logging is enabled for the Blob service.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.ChangeFeed.Enabled"),
			},
			{
				Name:        "blob_container_soft_delete_enabled",
				Description: "Specifies whether DeleteRetentionPolicy is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.ContainerDeleteRetentionPolicy.Enabled"),
			},
			{
				Name:        "blob_container_soft_delete_retention_days",
				Description: "Indicates the number of days that the deleted item should be retained.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.ContainerDeleteRetentionPolicy.Days"),
			},
			{
				Name:        "blob_restore_policy_days",
				Description: "Specifies how long the blob can be restored.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.RestorePolicy.Days"),
			},
			{
				Name:        "blob_restore_policy_enabled",
				Description: "Specifies whether blob restore is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.RestorePolicy.Enabled"),
			},
			{
				Name:        "blob_service_logging",
				Description: "Specifies the blob service properties for logging access.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAzureStorageAccountBlobServiceLogging,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "blob_soft_delete_enabled",
				Description: "Specifies whether DeleteRetentionPolicy is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.DeleteRetentionPolicy.Enabled"),
			},
			{
				Name:        "blob_soft_delete_retention_days",
				Description: "Indicates the number of days that the deleted item should be retained.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.DeleteRetentionPolicy.Days"),
			},
			{
				Name:        "blob_versioning_enabled",
				Description: "Specifies whether versioning is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountBlobProperties,
				Transform:   transform.FromField("BlobServicePropertiesProperties.IsVersioningEnabled"),
			},
			{
				Name:        "enable_https_traffic_only",
				Description: "Allows https traffic only to storage service if sets to true.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Account.AccountProperties.EnableHTTPSTrafficOnly"),
			},
			{
				Name:        "encryption_key_source",
				Description: "Contains the encryption keySource (provider).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.KeySource").Transform(transform.ToString),
			},
			{
				Name:        "encryption_key_vault_properties_key_current_version_id",
				Description: "The object identifier of the current versioned Key Vault Key in use.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.KeyVaultProperties.CurrentVersionedKeyIdentifier"),
			},
			{
				Name:        "encryption_key_vault_properties_key_name",
				Description: "The name of KeyVault key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.KeyVaultProperties.KeyName"),
			},
			{
				Name:        "encryption_key_vault_properties_key_vault_uri",
				Description: "The Uri of KeyVault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.KeyVaultProperties.KeyVaultURI"),
			},
			{
				Name:        "encryption_key_vault_properties_key_version",
				Description: "The version of KeyVault key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.KeyVaultProperties.KeyVersion"),
			},
			{
				Name:        "encryption_key_vault_properties_last_rotation_time",
				Description: "Timestamp of last rotation of the Key Vault Key.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.KeyVaultProperties.LastKeyRotationTimestamp").Transform(convertDateToTime),
			},
			{
				Name:        "failover_in_progress",
				Description: "Specifies whether the failover is in progress.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Account.AccountProperties.FailoverInProgress"),
			},
			{
				Name:        "file_soft_delete_enabled",
				Description: "Specifies whether DeleteRetentionPolicy is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountFileProperties,
				Transform:   transform.FromField("ShareDeleteRetentionPolicy.Enabled"),
			},
			{
				Name:        "file_soft_delete_retention_days",
				Description: "Indicates the number of days that the deleted item should be retained.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAzureStorageAccountFileProperties,
				Transform:   transform.FromField("ShareDeleteRetentionPolicy.Days"),
			},
			{
				Name:        "is_hns_enabled",
				Description: "Specifies whether account HierarchicalNamespace is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Account.AccountProperties.IsHnsEnabled"),
			},
			{
				Name:        "queue_logging_delete",
				Description: "Specifies whether all delete requests should be logged.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountQueueProperties,
				Transform:   transform.FromField("Logging.Delete"),
			},
			{
				Name:        "queue_logging_read",
				Description: "Specifies whether all read requests should be logged.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountQueueProperties,
				Transform:   transform.FromField("Logging.Read"),
			},
			{
				Name:        "queue_logging_retention_days",
				Description: "Indicates the number of days that metrics or logging data should be retained.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAzureStorageAccountQueueProperties,
				Transform:   transform.FromField("Logging.RetentionPolicy.Days"),
			},
			{
				Name:        "queue_logging_retention_enabled",
				Description: "Specifies whether a retention policy is enabled for the storage service.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountQueueProperties,
				Transform:   transform.FromField("Logging.RetentionPolicy.Enabled"),
			},
			{
				Name:        "queue_logging_version",
				Description: "The version of Storage Analytics to configure.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAzureStorageAccountQueueProperties,
				Transform:   transform.FromField("Logging.Version"),
			},
			{
				Name:        "queue_logging_write",
				Description: "Specifies whether all write requests should be logged.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAzureStorageAccountQueueProperties,
				Transform:   transform.FromField("Logging.Write"),
			},
			{
				Name:        "minimum_tls_version",
				Description: "Contains the minimum TLS version to be permitted on requests to storage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.MinimumTLSVersion").Transform(transform.ToString),
			},
			{
				Name:        "network_rule_bypass",
				Description: "Specifies whether traffic is bypassed for Logging/Metrics/AzureServices.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.NetworkRuleSet.Bypass").Transform(transform.ToString),
			},
			{
				Name:        "network_rule_default_action",
				Description: "Specifies the default action of allow or deny when no other rules match.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.NetworkRuleSet.DefaultAction").Transform(transform.ToString),
			},
			{
				Name:        "primary_blob_endpoint",
				Description: "Contains the blob endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.PrimaryEndpoints.Blob"),
			},
			{
				Name:        "primary_dfs_endpoint",
				Description: "Contains the dfs endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.PrimaryEndpoints.Dfs"),
			},
			{
				Name:        "primary_file_endpoint",
				Description: "Contains the file endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.PrimaryEndpoints.File"),
			},
			{
				Name:        "primary_location",
				Description: "Contains the location of the primary data center for the storage account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.PrimaryLocation"),
			},
			{
				Name:        "primary_queue_endpoint",
				Description: "Contains the queue endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.PrimaryEndpoints.Queue"),
			},
			{
				Name:        "primary_table_endpoint",
				Description: "Contains the table endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.PrimaryEndpoints.Table"),
			},
			{
				Name:        "primary_web_endpoint",
				Description: "Contains the web endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.PrimaryEndpoints.Web"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the virtual network resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "require_infrastructure_encryption",
				Description: "Specifies whether or not the service applies a secondary layer of encryption with platform managed keys for data at rest.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.RequireInfrastructureEncryption"),
			},
			{
				Name:        "secondary_location",
				Description: "Contains the location of the geo-replicated secondary for the storage account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.AccountProperties.SecondaryLocation"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the storage account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listStorageAccountDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption_scope",
				Description: "Encryption scope details for the storage account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAzureStorageAccountEncryptionScope,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption_services",
				Description: "A list of services which support encryption.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.AccountProperties.Encryption.Services"),
			},
			{
				Name:        "lifecycle_management_policy",
				Description: "The managementpolicy associated with the specified storage account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAzureStorageAccountLifecycleManagementPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "network_ip_rules",
				Description: "A list of IP ACL rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.AccountProperties.NetworkRuleSet.IPRules"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "A list of private endpoint connection associated with the specified storage account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.AccountProperties.PrivateEndpointConnections"),
			},
			{
				Name:        "virtual_network_rules",
				Description: "A list of virtual network rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.AccountProperties.NetworkRuleSet.VirtualNetworkRules"),
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
				Transform:   transform.FromField("Account.Tags"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Location").Transform(toLower),
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

//// LIST FUNCTION

func listStorageAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	logger := plugin.Logger(ctx)
	if err != nil {
		logger.Error("listStorageAccounts", "get session error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	storageClient := storage.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	result, err := storageClient.List(ctx)
	if err != nil {
		logger.Error("listStorageAccounts", "api error", err)
		return nil, err
	}

	for _, account := range result.Values() {
		resourceGroup := &strings.Split(string(*account.ID), "/")[4]
		d.StreamListItem(ctx, &storageAccountInfo{account, account.Name, resourceGroup})
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

		for _, account := range result.Values() {
			resourceGroup := &strings.Split(string(*account.ID), "/")[4]
			d.StreamListItem(ctx, &storageAccountInfo{account, account.Name, resourceGroup})
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

func getStorageAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageAccount")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	storageClient := storage.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	op, err := storageClient.GetProperties(ctx, resourceGroup, name, storage.AccountExpand("blobRestoreStatus"))
	if err != nil {
		return nil, err
	}

	return &storageAccountInfo{op, op.Name, &resourceGroup}, nil
}

func getAzureStorageAccountLifecycleManagementPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(*storageAccountInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewManagementPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	op, err := storageClient.Get(ctx, *accountData.ResourceGroup, *accountData.Name)
	if err != nil {
		if strings.Contains(err.Error(), "ManagementPolicyNotFound") {
			return nil, nil
		}
		return nil, err
	}

	// Direct assignment returns ManagementPolicyProperties only
	objectMap := make(map[string]interface{})
	if op.ID != nil {
		objectMap["id"] = op.ID
	}
	if op.Name != nil {
		objectMap["name"] = op.Name
	}
	if op.Type != nil {
		objectMap["type"] = op.Type
	}
	if op.ManagementPolicyProperties != nil {
		objectMap["properties"] = op.ManagementPolicyProperties
	}

	return objectMap, nil
}

func getAzureStorageAccountBlobProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(*storageAccountInfo)

	// Blob is not supported for the account if storage type is FileStorage
	if accountData.Account.Kind == "FileStorage" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewBlobServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	op, err := storageClient.GetServiceProperties(ctx, *accountData.ResourceGroup, *accountData.Name)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func listAzureStorageAccountEncryptionScope(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(*storageAccountInfo)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewEncryptionScopesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	encryptionScope, err := storageClient.List(ctx, *accountData.ResourceGroup, *accountData.Name)
	if err != nil {
		plugin.Logger(ctx).Error("listAzureStorageAccountEncryptionScope", "List", err)
		return nil, err
	}

	var encryptionScopes []map[string]interface{}

	for _, scope := range encryptionScope.Values() {
		encryptionScopes = append(encryptionScopes, storageAccountEncryptionScopeMap(scope))
	}

	for encryptionScope.NotDone() {
		err = encryptionScope.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listAzureStorageAccountEncryptionScope", "List_paging", err)
			return nil, err
		}
		for _, scope := range encryptionScope.Values() {
			encryptionScopes = append(encryptionScopes, storageAccountEncryptionScopeMap(scope))
		}
	}

	return encryptionScopes, nil
}

func getAzureStorageAccountBlobServiceLogging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(*storageAccountInfo)

	// Blob is not supported for the account if storage type is FileStorage
	if accountData.Account.Kind == "FileStorage" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	accountKeys, err := storageClient.ListKeys(ctx, *accountData.ResourceGroup, *accountData.Name, "")
	if err != nil {
		// storage.AccountsClient#ListKeys: Failure sending request: StatusCode=409 -- Original Error: autorest/azure: Service returned an error. Status=<nil> Code="ScopeLocked"
		// Message="The scope '/subscriptions/********-****-****-****-************/resourceGroups/turbot_rg/providers/Microsoft.Storage/storageAccounts/delmett'
		// cannot perform write operation because following scope(s) are locked: '/subscriptions/********-****-****-****-************/resourcegroups/turbot_rg/providers/Microsoft.Storage/storageAccounts/delmett'.
		// Please remove the lock and try again."
		if strings.Contains(err.Error(), "ScopeLocked") {
			return nil, nil
		}
		return nil, err
	}

	if *accountKeys.Keys != nil || len(*accountKeys.Keys) > 0 {
		key := (*accountKeys.Keys)[0]
		storageAuth, err := autorest.NewSharedKeyAuthorizer(*accountData.Name, *key.Value, autorest.SharedKeyLite)
		if err != nil {
			return nil, err
		}

		client := accounts.New()
		client.Client.Authorizer = storageAuth
		client.BaseURI = session.StorageEndpointSuffix

		resp, err := client.GetServiceProperties(ctx, *accountData.Name)
		if err != nil {
			if strings.Contains(err.Error(), "FeatureNotSupportedForAccount") {
				return nil, nil
			}
			return nil, err
		}
		return resp.StorageServiceProperties.Logging, nil
	}
	return nil, nil
}

func getAzureStorageAccountFileProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(*storageAccountInfo)

	// ge.FileServicesClient#GetServiceProperties: Failure responding to request: StatusCode=400 --
	// Original Error: autorest/azure: Service returned an error. Status=400 Code="FeatureNotSupportedForAccount" Message="File is not supported for the account."
	if accountData.Account.Kind == "BlobStorage" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewFileServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	op, err := storageClient.GetServiceProperties(ctx, *accountData.ResourceGroup, *accountData.Name)
	if err != nil {
		if strings.Contains(err.Error(), "FeatureNotSupportedForAccount") {
			return nil, nil
		}
		return nil, err
	}

	return op.FileServicePropertiesProperties, nil
}

func getAzureStorageAccountQueueProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(*storageAccountInfo)

	// ge.FileServicesClient#GetServiceProperties: Failure responding to request: StatusCode=400 --
	// Original Error: autorest/azure: Service returned an error. Status=400 Code="FeatureNotSupportedForAccount" Message="File is not supported for the account."
	if accountData.Account.Sku.Tier == "Standard" && (accountData.Account.Kind == "Storage" || accountData.Account.Kind == "StorageV2") {
		// Create session
		session, err := GetNewSession(ctx, d, "MANAGEMENT")
		if err != nil {
			return nil, err
		}
		subscriptionID := session.SubscriptionID

		storageClient := storage.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
		storageClient.Authorizer = session.Authorizer

		accountKeys, err := storageClient.ListKeys(ctx, *accountData.ResourceGroup, *accountData.Name, "")
		if err != nil {
			// storage.AccountsClient#ListKeys: Failure sending request: StatusCode=409 -- Original Error: autorest/azure: Service returned an error. Status=<nil> Code="ScopeLocked"
			// Message="The scope '/subscriptions/********-****-****-****-************/resourceGroups/turbot_rg/providers/Microsoft.Storage/storageAccounts/delmett'
			// cannot perform write operation because following scope(s) are locked: '/subscriptions/********-****-****-****-************/resourcegroups/turbot_rg/providers/Microsoft.Storage/storageAccounts/delmett'.
			// Please remove the lock and try again."
			if strings.Contains(err.Error(), "ScopeLocked") {
				return nil, nil
			}
			return nil, err
		}

		if *accountKeys.Keys != nil || len(*accountKeys.Keys) > 0 {
			key := (*accountKeys.Keys)[0]
			storageAuth, err := autorest.NewSharedKeyAuthorizer(*accountData.Name, *key.Value, autorest.SharedKeyLite)
			if err != nil {
				return nil, err
			}

			queuesClient := queues.New()
			queuesClient.Client.Authorizer = storageAuth
			queuesClient.BaseURI = session.StorageEndpointSuffix

			// using 	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/queues" to logging details
			// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage#QueueServicePropertiesProperties
			// In azure SDK for GO, we still don't have logging properties in its output
			resp, err := queuesClient.GetServiceProperties(ctx, *accountData.Name)

			if err != nil {
				if strings.Contains(err.Error(), "FeatureNotSupportedForAccount") {
					return nil, nil
				}
				return nil, err
			}
			return resp.StorageServiceProperties, nil
		}
	}
	return nil, nil
}

func listStorageAccountDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listStorageAccountDiagnosticSettings")
	accountData := h.Item.(*storageAccountInfo)
	id := *accountData.Account.ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives top level
	// contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["ID"] = i.ID
		}
		if i.Name != nil {
			objectMap["Name"] = i.Name
		}
		if i.Type != nil {
			objectMap["Type"] = i.Type
		}
		if i.DiagnosticSettings != nil {
			objectMap["DiagnosticSettings"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}

	return diagnosticSettings, nil
}

// If we return the API response directly, the output only gives the top level property
func storageAccountEncryptionScopeMap(scope storage.EncryptionScope) map[string]interface{} {
	objMap := make(map[string]interface{})
	if scope.ID != nil {
		objMap["Id"] = scope.ID
	}
	if scope.Name != nil {
		objMap["Name"] = scope.Name
	}
	if scope.Type != nil {
		objMap["Type"] = scope.Type
	}
	if scope.EncryptionScopeProperties != nil {
		if scope.EncryptionScopeProperties.Source != "" {
			objMap["Source"] = scope.EncryptionScopeProperties.Source
		}
		if scope.EncryptionScopeProperties.State != "" {
			objMap["State"] = scope.EncryptionScopeProperties.State
		}
		if scope.EncryptionScopeProperties.CreationTime != nil {
			objMap["CreationTime"] = scope.EncryptionScopeProperties.CreationTime
		}
		if scope.EncryptionScopeProperties.LastModifiedTime != nil {
			objMap["LastModifiedTime"] = scope.EncryptionScopeProperties.LastModifiedTime
		}
		if scope.EncryptionScopeProperties.KeyVaultProperties != nil {
			if scope.EncryptionScopeProperties.KeyVaultProperties.KeyURI != nil {
				objMap["KeyURI"] = scope.EncryptionScopeProperties.KeyVaultProperties.KeyURI
			}
		}
	}
	return objMap
}
