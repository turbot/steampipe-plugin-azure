package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-01-01/backup"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureRecoveryServicesVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_recovery_services_vault",
		Description: "Azure Recovery Services Vault",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getRecoveryServicesVault,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listRecoveryServicesVaults,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the recovery services vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "private_endpoint_state_for_site_recovery",
				Description: "Private endpoint state for site recovery of the recovery services vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PrivateEndpointStateForSiteRecovery"),
			},
			{
				Name:        "private_endpoint_state_for_backup",
				Description: "Private endpoint state for backup of the recovery services vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PrivateEndpointStateForBackup"),
			},
			{
				Name:        "sku_name",
				Description: "The sku name of the recovery services vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "backup_jobs",
				Description: "Backup jobs of the recovery services vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listRecoveryServicesVaultBackupJobs,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the recovery services vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listRecoveryServicesVaultDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity",
				Description: "Managed service identity of the recovery services vault.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_endpoint_connections",
				Description: "List of private endpoint connections of the recovery services vault.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PrivateEndpointConnections"),
			},
			{
				Name:        "upgrade_details",
				Description: "Upgrade details properties of the recovery services vault.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.UpgradeDetails"),
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

			// Azure standard column
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

//// LIST FUNCTION

func listRecoveryServicesVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	recoveryServicesVaultClient := recoveryservices.NewVaultsClient(subscriptionID)
	recoveryServicesVaultClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := recoveryServicesVaultClient.ListBySubscriptionID(context.Background())
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

//// HYDRATE FUNCTIONS

func getRecoveryServicesVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRecoveryServicesVault")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	recoveryServicesVaultClient := recoveryservices.NewVaultsClient(subscriptionID)
	recoveryServicesVaultClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := recoveryServicesVaultClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listRecoveryServicesVaultBackupJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vault := h.Item.(recoveryservices.Vault)
	resourceGroup := strings.Split(*vault.ID, "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	backupJobClient := backup.NewJobsClient(subscriptionID)
	backupJobClient.Authorizer = session.Authorizer

	// If we return the API response directly, the output only gives
	// the contents of BackupJobs
	var backupJobs []map[string]interface{}
	pagesLeft := true
	for pagesLeft {
		result, err := backupJobClient.List(ctx, *vault.Name, resourceGroup, "", "")
		if err != nil {
			return nil, err
		}
		for _, vault := range result.Values() {
			backupJob := make(map[string]interface{})
			if vault.ID != nil {
				backupJob["id"] = vault.ID
			}
			if vault.Name != nil {
				backupJob["name"] = vault.Name
			}
			if vault.Type != nil {
				backupJob["type"] = vault.Type
			}
			if vault.Location != nil {
				backupJob["Location"] = vault.Location
			}
			if vault.Tags != nil {
				backupJob["Tags"] = vault.Tags
			}
			if vault.ETag != nil {
				backupJob["ETag"] = vault.ETag
			}
			// if vault.Properties != nil {
			// 	backupJob["properties"] = vault.Properties
			// }
			backupJobs = append(backupJobs, backupJob)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return backupJobs, nil
}

func listRecoveryServicesVaultDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRecoveryServicesVaultDiagnosticSettings")
	id := *h.Item.(recoveryservices.Vault).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
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
