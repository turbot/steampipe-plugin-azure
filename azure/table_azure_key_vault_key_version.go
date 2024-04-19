package azure

import (
	"context"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/mgmt/keyvault"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureKeyVaultKeyVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault_key_version",
		Description: "Azure Key Vault Key Version",
		List: &plugin.ListConfig{
			Hydrate:       listKeyVaultKeyVersions,
			ParentHydrate: listKeyVaults,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name: "key_name", Require: plugin.Optional,
				},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the key version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_name",
				Description: "The friendly name that identifies the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractKeyNameAndKeyIdFromVersionID, "KeyName"),
			},
			{
				Name:        "key_id",
				Description: "Contains ID to identify a key uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractKeyNameAndKeyIdFromVersionID, "KeyId"),
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a key version uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "vault_name",
				Description: "The friendly name that identifies the vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractVaultNameFromID),
			},
			{
				Name:        "enabled",
				Description: "Indicates whether the key version is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("KeyProperties.Attributes.Enabled"),
			},
			{
				Name:        "key_type",
				Description: "The type of the key. Possible values are: 'EC', 'ECHSM', 'RSA', 'RSAHSM'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultKeyVersion,
				Transform:   transform.FromField("KeyProperties.Kty").Transform(transform.ToString),
			},
			{
				Name:        "created_at",
				Description: "Specifies the time when the key version is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.Created").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "curve_name",
				Description: "The elliptic curve name. Possible values are: 'P256', 'P384', 'P521', 'P256K'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVaultKeyVersion,
				Transform:   transform.FromField("KeyProperties.CurveName").Transform(transform.ToString),
			},
			{
				Name:        "expires_at",
				Description: "Specifies the time when the key version wil expire.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.Expires").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "key_size",
				Description: "The key size in bits.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKeyVaultKeyVersion,
				Transform:   transform.FromField("KeyProperties.KeySize"),
			},
			{
				Name:        "key_uri",
				Description: "The URI to retrieve the current version of the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyProperties.KeyURI"),
			},
			{
				Name:        "key_uri_with_version",
				Description: "The URI to retrieve the specific version of the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyProperties.KeyURIWithVersion"),
			},
			{
				Name:        "location",
				Description: "Azure location of the key vault resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "not_before",
				Description: "Specifies the time before which the key version is not usable.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.NotBefore").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "recovery_level",
				Description: "The deletion recovery level currently in effect for the object. If it contains 'Purgeable', then the object can be permanently deleted by a privileged user; otherwise, only the system can purge the object at the end of the retention interval.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyProperties.Attributes.RecoveryLevel").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "Specifies the time when the key was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("KeyProperties.Attributes.Updated").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "key_ops",
				Description: "A list of key operations.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVaultKeyVersion,
				Transform:   transform.FromField("KeyProperties.KeyOps"),
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

//// LIST FUNCTION

func listKeyVaultKeyVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of key
	vault := h.Item.(keyvault.Resource)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_key_version.listKeyVaultKeyVersions", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*vault.ID, "/")[4]

	client := keyvault.NewKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	var keys []keyvault.Key
	result, err := client.List(ctx, resourceGroup, *vault.Name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_key_version.listKeyVaultKeyVersions", "api_error", err)
		return nil, err
	}
	keys = append(keys, result.Values()...)

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		keys = append(keys, result.Values()...)
	}

	var wg sync.WaitGroup
	keyVersionCh := make(chan []keyvault.Key, len(keys))
	errorCh := make(chan error, len(keys))

	// Iterating all the available keys
	for _, item := range keys {
		wg.Add(1)
		go getRowDataForKeyVersionAsync(ctx, d, h, item, &wg, keyVersionCh, errorCh)
	}

	// wait for all executions to be processed
	wg.Wait()
	close(keyVersionCh)
	close(errorCh)

	for err := range errorCh {
		return nil, err
	}

	for item := range keyVersionCh {
		for _, data := range item {
			d.StreamLeafListItem(ctx, data)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

func getRowDataForKeyVersionAsync(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, key keyvault.Key, wg *sync.WaitGroup, keyVersionCh chan []keyvault.Key, errorCh chan error) {
	defer wg.Done()

	rowData, err := getRowDataForKeyVersion(ctx, d, h, key)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		keyVersionCh <- rowData
	}
}

func getRowDataForKeyVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, key keyvault.Key) ([]keyvault.Key, error) {
	vault := h.Item.(keyvault.Resource)
	keyName := d.EqualsQuals["key_name"].GetStringValue()
	var items []keyvault.Key

	if keyName != "" && keyName != *key.Name {
		return items, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_key_version.getRowDataForKeyVersion", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := strings.Split(*vault.ID, "/")[4]

	client := keyvault.NewKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListVersions(ctx, resourceGroup, *vault.Name, *key.Name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_key_version.getRowDataForKeyVersion", "api_error", err)
		return nil, err
	}

	items = append(items, op.Values()...)
	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		items = append(items, op.Values()...)
	}

	return items, nil
}

//// HYDRATE FUNCTIONS

func getKeyVaultKeyVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var vaultName, name, resourceGroup, keyVersion string
	data := h.Item.(keyvault.Key)
	splitID := strings.Split(*data.ID, "/")
	vaultName = splitID[8]
	name = splitID[10]
	resourceGroup = splitID[4]
	keyVersion = *data.Name

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_key_version.getKeyVaultKeyVersion", "client_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := keyvault.NewKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.GetVersion(ctx, resourceGroup, vaultName, name, keyVersion)
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault_key_version.getKeyVaultKeyVersion", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func extractKeyNameAndKeyIdFromVersionID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	data := d.HydrateItem.(keyvault.Key)
	splitVersion := strings.Split(*data.ID, "/")
	removedVersion := splitVersion[:len(splitVersion)-2]
	keyName := strings.Split(*data.ID, "/")[10]
	keyInfo := map[string]string{
		"KeyName": keyName,
		"KeyId":   strings.Join(removedVersion[:], "/"),
	}

	return keyInfo[param], nil
}
