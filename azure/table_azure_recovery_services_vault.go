package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/recoveryservices"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureRecoveryServicesVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_recovery_services_vault",
		Description: "Azure Recovery Services Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getRecoveryServicesVault,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listRecoveryServicesVaults,
		},
		Columns: azureColumns([]*plugin.Column{
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

func listRecoveryServicesVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	recoveryServicesVaultClient := recoveryservices.NewVaultsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	recoveryServicesVaultClient.Authorizer = session.Authorizer
	result, err := recoveryServicesVaultClient.ListBySubscriptionID(ctx)
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

func getRecoveryServicesVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRecoveryServicesVault")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	recoveryServicesVaultClient := recoveryservices.NewVaultsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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

func listRecoveryServicesVaultDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRecoveryServicesVaultDiagnosticSettings")
	id := *h.Item.(recoveryservices.Vault).ID

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
