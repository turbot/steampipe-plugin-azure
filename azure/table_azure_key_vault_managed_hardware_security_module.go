package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureKeyVaultManagedHardwareSecurityModule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault_managed_hardware_security_module",
		Description: "Azure Key Vault Managed Hardware Security Module",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getKeyVaultManagedHardwareSecurityModule,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyVaultManagedHardwareSecurityModules,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the managed HSM Pool.",
			},
			{
				Name:        "id",
				Description: "The Azure Resource Manager resource ID for the managed HSM Pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the managed HSM Pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state. Possible values include: 'ProvisioningStateSucceeded', 'ProvisioningStateProvisioning', 'ProvisioningStateFailed', 'ProvisioningStateUpdating', 'ProvisioningStateDeleting', 'ProvisioningStateActivated', 'ProvisioningStateSecurityDomainRestore', 'ProvisioningStateRestoring'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "hsm_uri",
				Description: "The URI of the managed hsm pool for performing operations on keys.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HsmURI"),
			},
			{
				Name:        "enable_soft_delete",
				Description: "Property to specify whether the 'soft delete' functionality is enabled for this managed HSM pool. If it's not set to any value(true or false) when creating new managed HSM pool, it will be set to true by default. Once set to true, it cannot be reverted to false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnableSoftDelete"),
			},
			{
				Name:        "soft_delete_retention_in_days",
				Description: "Indicates softDelete data retention days. It accepts >=7 and <=90.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.SoftDeleteRetentionInDays"),
			},
			{
				Name:        "enable_purge_protection",
				Description: "Property specifying whether protection against purge is enabled for this managed HSM pool. Setting this property to true activates protection against purge for this managed HSM pool and its content - only the Managed HSM service may initiate a hard, irrecoverable deletion. The setting is effective only if soft delete is also enabled. Enabling this functionality is irreversible.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnablePurgeProtection"),
			},
			{
				Name:        "status_message",
				Description: "Resource Status Message.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.StatusMessage"),
			},
			{
				Name:        "create_mode",
				Description: "The create mode to indicate whether the resource is being created or is being recovered from a deleted resource. Possible values include: 'CreateModeRecover', 'CreateModeDefault'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CreateMode"),
			},
			{
				Name:        "sku_family",
				Description: "Contains SKU family name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Sku.Family"),
			},
			{
				Name:        "sku_name",
				Description: "SKU name to specify whether the key vault is a standard vault or a premium vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "tenant_id",
				Description: "The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.TenantID").Transform(transform.ToString),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the managed HSM.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listKeyVaultHsmDiagnosticSettings,
				Transform:   transform.FromValue(),
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

func listKeyVaultManagedHardwareSecurityModules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	hsmClient := keyvault.NewManagedHsmsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	hsmClient.Authorizer = session.Authorizer
	maxResults := int32(100)

	result, err := hsmClient.ListBySubscription(ctx, &maxResults)
	if err != nil {
		return nil, err
	}
	for _, vault := range result.Values() {
		d.StreamListItem(ctx, vault)
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
		for _, vault := range result.Values() {
			d.StreamListItem(ctx, vault)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKeyVaultManagedHardwareSecurityModule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVaultManagedHardwareSecurityModule")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	var name, resourceGroup string
	if h.Item != nil {
		data := h.Item.(keyvault.ManagedHsm)
		name = *data.Name
		resourceGroup = strings.Split(*data.ID, "/")[4]
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroup = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	client := keyvault.NewManagedHsmsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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

func listKeyVaultHsmDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKmsKeyVaultHsmDiagnosticSettings")
	id := h.Item.(keyvault.ManagedHsm).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, *id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
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
		if i.DiagnosticSettings != nil {
			objectMap["properties"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}
	return diagnosticSettings, nil
}
