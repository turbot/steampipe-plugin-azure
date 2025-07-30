package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/mgmt/keyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureKeyVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault",
		Description: "Azure Key Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getKeyVault,
			Tags: map[string]string{
				"service": "Microsoft.KeyVault",
				"action":  "vaults/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyVaults,
			Tags: map[string]string{
				"service": "Microsoft.KeyVault",
				"action":  "vaults/read",
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listKmsKeyVaultDiagnosticSettings,
				Tags: map[string]string{
					"service": "Microsoft.Insights",
					"action":  "diagnosticSettings/read",
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the vault.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a vault uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "vault_uri",
				Description: "Contains URI of the vault for performing operations on keys and secrets.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.VaultURI"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_mode",
				Description: "The vault's create mode to indicate whether the vault need to be recovered or not. Possible values include: 'default', 'recover'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.CreateMode"),
			},
			{
				Name:        "enabled_for_deployment",
				Description: "Indicates whether Azure Virtual Machines are permitted to retrieve certificates stored as secrets from the key vault.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.EnabledForDeployment"),
				Default:     false,
			},
			{
				Name:        "enabled_for_disk_encryption",
				Description: "Indicates whether Azure Disk Encryption is permitted to retrieve secrets from the vault and unwrap keys.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.EnabledForDiskEncryption"),
				Default:     false,
			},
			{
				Name:        "enabled_for_template_deployment",
				Description: "Indicates whether Azure Resource Manager is permitted to retrieve secrets from the key vault.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.EnabledForTemplateDeployment"),
				Default:     false,
			},
			{
				Name:        "enable_rbac_authorization",
				Description: "Property that controls how data actions are authorized.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.EnableRbacAuthorization"),
				Default:     false,
			},
			{
				Name:        "purge_protection_enabled",
				Description: "Indicates whether protection against purge is enabled for this vault.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.EnablePurgeProtection"),
				Default:     false,
			},
			{
				Name:        "soft_delete_enabled",
				Description: "Indicates whether the 'soft delete' functionality is enabled for this key vault.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.EnableSoftDelete"),
			},
			{
				Name:        "soft_delete_retention_in_days",
				Description: "Contains softDelete data retention days.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.SoftDeleteRetentionInDays"),
			},
			{
				Name:        "sku_family",
				Description: "Contains SKU family name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.Sku.Family"),
			},
			{
				Name:        "sku_name",
				Description: "SKU name to specify whether the key vault is a standard vault or a premium vault.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "tenant_id",
				Description: "The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.TenantID").Transform(transform.ToString),
			},
			{
				Name:        "access_policies",
				Description: "A list of 0 to 1024 identities that have access to the key vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVault,
				Transform:   transform.From(extractKeyVaultAccessPolicies),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listKmsKeyVaultDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "network_acls",
				Description: "Rules governing the accessibility of the key vault from specific network locations.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVault,
				Transform:   transform.FromField("Properties.NetworkAcls"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "List of private endpoint connections associated with the key vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyVault,
				Transform:   transform.From(extractKeyVaultPrivateEndpointConnections),
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

type PrivateEndpointConnectionInfo struct {
	PrivateEndpointId                               string
	PrivateLinkServiceConnectionStateStatus         string
	PrivateLinkServiceConnectionStateDescription    string
	PrivateLinkServiceConnectionStateActionRequired string
	ProvisioningState                               string
}

//// LIST FUNCTION

func listKeyVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	keyVaultClient := keyvault.NewVaultsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	keyVaultClient.Authorizer = session.Authorizer
	maxResults := int32(100)

	// Apply Retry rule
	ApplyRetryRules(ctx, &keyVaultClient, d.Connection)

	result, err := keyVaultClient.List(ctx, &maxResults)
	if err != nil {
		return nil, err
	}
	for _, vault := range result.Values() {
		d.StreamListItem(ctx, vault)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, vault := range result.Values() {
			d.StreamListItem(ctx, vault)
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

func getKeyVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVault")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	var name, resourceGroup string
	if h.Item != nil {
		data := h.Item.(keyvault.Resource)
		name = *data.Name
		resourceGroup = strings.Split(*data.ID, "/")[4]
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		resourceGroup = d.EqualsQuals["resource_group"].GetStringValue()
	}

	client := keyvault.NewVaultsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

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

func listKmsKeyVaultDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKmsKeyVaultDiagnosticSettings")
	id := getKeyVaultID(h.Item)

	// Create session
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armmonitor.NewDiagnosticSettingsClient(session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_key_vault.listKmsKeyVaultDiagnosticSettings", "client_error", err)
		return nil, err
	}

	var diagnosticSettings []map[string]interface{}

	input := &armmonitor.DiagnosticSettingsClientListOptions{}

	pager := clientFactory.NewListPager(id, input)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_key_vault.listKmsKeyVaultDiagnosticSettings", "api_error", err)
			return nil, err
		}
		for _, i := range page.Value {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Type != nil {
				objectMap["type"] = i.Type
			}
			if i.Properties != nil {
				objectMap["properties"] = i.Properties
			}
			diagnosticSettings = append(diagnosticSettings, objectMap)
		}
	}
	return diagnosticSettings, nil
}

//// TRANSFORM FUNCTIONS

func extractKeyVaultPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	vault := d.HydrateItem.(keyvault.Vault)
	plugin.Logger(ctx).Trace("extractKeyVaultPrivateEndpointConnections")
	var privateEndpointDetails []PrivateEndpointConnectionInfo
	var privateEndpoint PrivateEndpointConnectionInfo
	if vault.Properties.PrivateEndpointConnections != nil {
		for _, connection := range *vault.Properties.PrivateEndpointConnections {
			// Below checks are required for handling invalid memory address or nil pointer dereference error
			if connection.PrivateEndpointConnectionProperties != nil {
				if connection.PrivateEndpoint != nil {
					privateEndpoint.PrivateEndpointId = *connection.PrivateEndpoint.ID
				}
				if connection.PrivateLinkServiceConnectionState != nil {
					if connection.PrivateLinkServiceConnectionState.ActionsRequired != "" {
						privateEndpoint.PrivateLinkServiceConnectionStateActionRequired = string(connection.PrivateLinkServiceConnectionState.ActionsRequired)
					}
					if connection.PrivateLinkServiceConnectionState.Description != nil {
						privateEndpoint.PrivateLinkServiceConnectionStateDescription = *connection.PrivateLinkServiceConnectionState.Description
					}
					if connection.PrivateLinkServiceConnectionState.Status != "" {
						privateEndpoint.PrivateLinkServiceConnectionStateStatus = string(connection.PrivateLinkServiceConnectionState.Status)
					}
				}
				if connection.ProvisioningState != "" {
					privateEndpoint.ProvisioningState = string(connection.ProvisioningState)
				}
			}
			privateEndpointDetails = append(privateEndpointDetails, privateEndpoint)
		}
	}

	return privateEndpointDetails, nil
}

// If we return the API response directly, the output will not provide the properties of AccessPolicies
func extractKeyVaultAccessPolicies(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	vault := d.HydrateItem.(keyvault.Vault)
	var policies []map[string]interface{}

	if vault.Properties.AccessPolicies != nil {
		for _, i := range *vault.Properties.AccessPolicies {
			objectMap := make(map[string]interface{})
			if i.TenantID != nil {
				objectMap["tenantId"] = i.TenantID
			}
			if i.ObjectID != nil {
				objectMap["objectId"] = i.ObjectID
			}
			if i.ApplicationID != nil {
				objectMap["applicationId"] = i.ApplicationID
			}
			if i.Permissions != nil {
				if i.Permissions.Keys != nil {
					objectMap["permissionsKeys"] = i.Permissions.Keys
				}
				if i.Permissions.Secrets != nil {
					objectMap["permissionsSecrets"] = i.Permissions.Secrets
				}
				if i.Permissions.Keys != nil {
					objectMap["permissionsCertificates"] = i.Permissions.Certificates
				}
				if i.Permissions.Keys != nil {
					objectMap["permissionsStorage"] = i.Permissions.Storage
				}
			}
			policies = append(policies, objectMap)
		}
	}

	return policies, nil
}

func getKeyVaultID(item interface{}) string {
	switch item := item.(type) {
	case keyvault.Vault:
		return *item.ID
	case keyvault.Resource:
		return *item.ID
	}
	return ""
}