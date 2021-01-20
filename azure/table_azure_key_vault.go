package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureKeyVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_key_vault",
		Description: "Azure Key Vault",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getKeyVault,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyVaults,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the vault",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a vault uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "vault_uri",
				Description: "Contains URI of the vault for performing operations on keys and secrets",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.VaultURI"),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled_for_deployment",
				Description: "Indicates whether Azure Virtual Machines are permitted to retrieve certificates stored as secrets from the key vault",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnabledForDeployment"),
				Default:     false,
			},
			{
				Name:        "enabled_for_disk_encryption",
				Description: "Indicates whether Azure Disk Encryption is permitted to retrieve secrets from the vault and unwrap keys",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnabledForDiskEncryption"),
				Default:     false,
			},
			{
				Name:        "enabled_for_template_deployment",
				Description: "Indicates whether Azure Resource Manager is permitted to retrieve secrets from the key vault",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnabledForTemplateDeployment"),
				Default:     false,
			},
			{
				Name:        "enable_rbac_authorization",
				Description: "Property that controls how data actions are authorized",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnableRbacAuthorization"),
				Default:     false,
			},
			{
				Name:        "purge_protection_enabled",
				Description: "Indicates whether protection against purge is enabled for this vault",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnablePurgeProtection"),
				Default:     false,
			},
			{
				Name:        "soft_delete_enabled",
				Description: "Indicates whether the 'soft delete' functionality is enabled for this key vault",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnableSoftDelete"),
			},
			{
				Name:        "soft_delete_retention_in_days",
				Description: "Contains softDelete data retention days",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.SoftDeleteRetentionInDays"),
			},
			{
				Name:        "sku_family",
				Description: "Contains SKU family name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Sku.Family"),
			},
			{
				Name:        "sku_name",
				Description: "SKU name to specify whether the key vault is a standard vault or a premium vault",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "tenant_id",
				Description: "The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.TenantID").Transform(transform.ToString),
			},
			{
				Name:        "access_policies",
				Description: "A list of 0 to 1024 identities that have access to the key vault",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AccessPolicies"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},
			{
				Name:        "region",
				Description: "The Azure region in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "resource_group",
				Description: "Name of the resource group, the vault is created at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
			{
				Name:        "subscription_id",
				Description: "The Azure Subscription ID in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// FETCH FUNCTIONS ////

func listKeyVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	keyVaultClient := keyvault.NewVaultsClient(subscriptionID)
	keyVaultClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := keyVaultClient.ListBySubscription(ctx, nil)
		if err != nil {
			return nil, err
		}

		for _, vault := range result.Values() {
			d.StreamListItem(ctx, vault)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getKeyVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVault")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	keyVaultClient := keyvault.NewVaultsClient(subscriptionID)
	keyVaultClient.Authorizer = session.Authorizer

	op, err := keyVaultClient.Get(ctx, resourceGroup, name)
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
