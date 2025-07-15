package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cognitiveservices/mgmt/cognitiveservices"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureCognitiveAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cognitive_account",
		Description: "Azure Cognitive Account",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getCognitiveAccount,
			Tags: map[string]string{
				"service": "Microsoft.CognitiveServices",
				"action":  "accounts/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listCognitiveAccounts,
			Tags: map[string]string{
				"service": "Microsoft.CognitiveServices",
				"action":  "accounts/read",
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listCognitiveAccountDiagnosticSettings,
				Tags: map[string]string{
					"service": "Microsoft.Insights",
					"action":  "diagnosticSettings/read",
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "kind",
				Description: "The kind of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The status of the cognitive services account at the time the operation was called. Possible values include: 'Accepted', 'Creating', 'Deleting', 'Moving', 'Failed', 'Succeeded', 'ResolvingDNS'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource. E.g. 'Microsoft.Compute/virtualMachines' or 'Microsoft.Storage/storageAccounts'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "custom_sub_domain_name",
				Description: "The subdomain name used for token-based authentication.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CustomSubDomainName"),
			},
			{
				Name:        "date_created",
				Description: "The date of cognitive services account creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DateCreated"),
			},
			{
				Name:        "disable_local_auth",
				Description: "Checks if local auth is disabled for the resource.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Properties.DisableLocalAuth"),
			},
			{
				Name:        "endpoint",
				Description: "The endpoint of the created account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Endpoint"),
			},
			{
				Name:        "etag",
				Description: "The resource etag.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_migrated",
				Description: "Checks if the resource is migrated from an existing key.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Properties.IsMigrated"),
			},
			{
				Name:        "migration_token",
				Description: "The resource migration token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MigrationToken"),
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public endpoint access is allowed for this account. Value is optional but if passed in, must be 'Enabled' or 'Disabled'. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PublicNetworkAccess"),
			},
			{
				Name:        "restore",
				Description: "Checks if restore is enabled for the resource.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Properties.Restore"),
			},
			{
				Name:        "restrict_outbound_network_access",
				Description: "Checks if outbound network access is restricted for the resource.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Properties.RestrictOutboundNetworkAccess"),
			},
			{
				Name:        "allowed_fqdn_list",
				Description: "The allowed FQDN list for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AllowedFqdnList"),
			},
			{
				Name:        "api_properties",
				Description: "The api properties for special APIs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.APIProperties"),
			},
			{
				Name:        "call_rate_limit",
				Description: "The call rate limit of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CallRateLimit"),
			},
			{
				Name:        "capabilities",
				Description: "The capabilities of the cognitive services account. Each item indicates the capability of a specific feature. The values are read-only and for reference only.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Capabilities"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the cognitive service account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listCognitiveAccountDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "The encryption properties for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Encryption"),
			},
			{
				Name:        "endpoints",
				Description: "All endpoints of the cognitive services account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Endpoints"),
			},
			{
				Name:        "identity",
				Description: "The identity for the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_acls",
				Description: "A collection of rules governing the accessibility from specific network locations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.NetworkAcls"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The private endpoint connection associated with the cognitive services account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractAccountPrivateEndpointConnections),
			},
			{
				Name:        "quota_limit",
				Description: "The quota limit of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.QuotaLimit"),
			},
			{
				Name:        "sku",
				Description: "The resource model definition representing SKU.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sku_change_info",
				Description: "Sku change info of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.SkuChangeInfo"),
			},
			{
				Name:        "system_data",
				Description: "The metadata pertaining to creation and last modification of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_owned_storage",
				Description: "The storage accounts for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.UserOwnedStorage"),
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

type CognitiveAccountPrivateEndpointConnections struct {
	PrivateEndpointID                 interface{}
	PrivateLinkServiceConnectionState interface{}
	ProvisioningState                 interface{}
	GroupIds                          *[]string
	SystemData                        interface{}
	Location                          *string
	Etag                              *string
	ID                                *string
	Name                              *string
	Type                              *string
}

//// LIST FUNCTION

func listCognitiveAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountsClient := cognitiveservices.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountsClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &accountsClient, d.Connection)

	result, err := accountsClient.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listCognitiveAccounts", "list", err)
		return nil, err
	}

	for _, account := range result.Values() {
		d.StreamListItem(ctx, account)
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listCognitiveAccounts", "list_paging", err)
			return nil, err
		}
		for _, account := range result.Values() {
			d.StreamListItem(ctx, account)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCognitiveAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCognitiveAccount")

	accountName := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if accountName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountsClient := cognitiveservices.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountsClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &accountsClient, d.Connection)

	account, err := accountsClient.Get(ctx, resourceGroup, accountName)
	if err != nil {
		plugin.Logger(ctx).Error("getCognitiveAccount", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if account.ID != nil {
		return account, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func listCognitiveAccountDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listCognitiveAccountDiagnosticSettings")
	id := *h.Item.(cognitiveservices.Account).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	result, err := client.List(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("listCognitiveAccountDiagnosticSettings", "list", err)
		return nil, err
	}

	// If we return the API response directly, the output does not provide
	// the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *result.Value {
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

// If we return the API response directly, the output will not provide all the properties of PrivateEndpointConnections
func extractAccountPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	account := d.HydrateItem.(cognitiveservices.Account)
	privateEndpointConnectionInfo := []CognitiveAccountPrivateEndpointConnections{}

	if account.Properties.PrivateEndpointConnections != nil {
		for _, connection := range *account.Properties.PrivateEndpointConnections {
			properties := CognitiveAccountPrivateEndpointConnections{}
			properties.SystemData = connection.SystemData
			properties.Location = connection.Location
			properties.Etag = connection.Etag
			properties.ID = connection.ID
			properties.Name = connection.Name
			properties.Type = connection.Type
			if connection.Properties != nil {
				if connection.Properties.PrivateEndpoint != nil {
					properties.PrivateEndpointID = connection.Properties.PrivateEndpoint.ID
				}
				properties.PrivateLinkServiceConnectionState = connection.Properties.PrivateLinkServiceConnectionState
				properties.ProvisioningState = connection.Properties.ProvisioningState
				properties.GroupIds = connection.Properties.GroupIds
			}
			privateEndpointConnectionInfo = append(privateEndpointConnectionInfo, properties)
		}
	}

	return privateEndpointConnectionInfo, nil
}
