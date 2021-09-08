package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2021-04-30/cognitiveservices"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureCognitiveAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cognitive_account",
		Description: "Azure Cognitive Account",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getCognitiveAccount,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listCognitiveAccounts,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Name"),
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.ID"),
			},
			{
				Name:        "kind",
				Description: "The kind of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Kind"),
			},
			{
				Name:        "provisioning_state",
				Description: "Gets the status of the cognitive services account at the time the operation was called. Possible values include: 'Accepted', 'Creating', 'Deleting', 'Moving', 'Failed', 'Succeeded', 'ResolvingDNS'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource. E.g. 'Microsoft.Compute/virtualMachines' or 'Microsoft.Storage/storageAccounts'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Type"),
			},
			{
				Name:        "custom_sub_domain_name",
				Description: "The subdomain name used for token-based authentication.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Properties.CustomSubDomainName"),
			},
			{
				Name:        "date_created",
				Description: "Gets the date of cognitive services account creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Properties.DateCreated"),
			},
			{
				Name:        "disable_local_auth",
				Description: "Checks if local auth is disabled for the resource.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Account.Properties.DisableLocalAuth"),
			},
			{
				Name:        "endpoint",
				Description: "The endpoint of the created account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Properties.Endpoint"),
			},
			{
				Name:        "etag",
				Description: "The resource etag.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Etag"),
			},
			{
				Name:        "is_migrated",
				Description: "Checks if the resource is migrated from an existing key.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Account.Properties.IsMigrated"),
			},
			{
				Name:        "migration_token",
				Description: "The resource migration token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Properties.MigrationToken"),
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public endpoint access is allowed for this account. Value is optional but if passed in, must be 'Enabled' or 'Disabled'. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Properties.PublicNetworkAccess"),
			},
			{
				Name:        "restore",
				Description: "Checks if restore is enabled for the resource.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Account.Properties.Restore"),
			},
			{
				Name:        "restrict_outbound_network_access",
				Description: "Checks if outbound network access is restricted for the resource.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Account.Properties.RestrictOutboundNetworkAccess"),
			},
			{
				Name:        "allowed_fqdn_list",
				Description: "The allowed FQDN list for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.AllowedFqdnList"),
			},
			{
				Name:        "api_properties",
				Description: "The api properties for special APIs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.APIProperties"),
			},
			{
				Name:        "call_rate_limit",
				Description: "The call rate limit of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.CallRateLimit"),
			},
			{
				Name:        "capabilities",
				Description: "Gets the capabilities of the cognitive services account. Each item indicates the capability of a specific feature. The values are read-only and for reference only.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.Capabilities"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listCognitiveAccountDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "The encryption properties for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.Encryption"),
			},
			{
				Name:        "endpoints",
				Description: "All endpoints of the cognitive services account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.Endpoints"),
			},
			{
				Name:        "identity",
				Description: "The identity for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Identity"),
			},
			{
				Name:        "network_acls",
				Description: "A collection of rules governing the accessibility from specific network locations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.NetworkAcls"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The private endpoint connection associated with the cognitive services account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointConnections"),
			},
			{
				Name:        "quota_limit",
				Description: "The quota limit of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.QuotaLimit"),
			},
			{
				Name:        "sku",
				Description: "The resource model definition representing SKU.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Sku"),
			},
			{
				Name:        "sku_change_info",
				Description: "Sku change info of account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.SkuChangeInfo"),
			},
			{
				Name:        "system_data",
				Description: "The metadata pertaining to creation and last modification of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.SystemData"),
			},
			{
				Name:        "user_owned_storage",
				Description: "The storage accounts for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Account.Properties.UserOwnedStorage"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.Name"),
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
				Transform:   transform.FromField("Account.ID").Transform(extractResourceGroupFromID),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

type AccountInfo struct {
	Account                    cognitiveservices.Account
	PrivateEndpointConnections []PrivateEndpointConnectionInfo
}

type PrivateEndpointConnectionInfo struct {
	Properties PrivateEndpointConnectionPropertiesInfo
	SystemData interface{}
	Location   *string
	Etag       *string
	ID         *string
	Name       *string
	Type       *string
}

type PrivateEndpointConnectionPropertiesInfo struct {
	PrivateEndpoint                   PrivateEndpointInfo
	PrivateLinkServiceConnectionState interface{}
	ProvisioningState                 interface{}
	GroupIds                          *[]string
}

type PrivateEndpointInfo struct {
	ID *string
}

//// LIST FUNCTIONS

func listCognitiveAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountsClient := cognitiveservices.NewAccountsClient(subscriptionID)
	accountsClient.Authorizer = session.Authorizer

	result, err := accountsClient.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, account := range result.Values() {
		privateEndpointConnectionInfo := getPrivateEndpointConnectionInfo(&account)
		d.StreamListItem(ctx, AccountInfo{account, privateEndpointConnectionInfo})
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, account := range result.Values() {
			privateEndpointConnectionInfo := getPrivateEndpointConnectionInfo(&account)
			d.StreamListItem(ctx, AccountInfo{account, privateEndpointConnectionInfo})
		}
	}
	
	return nil, err
}

//// HYDRATE FUNCTIONS

func getCognitiveAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCognitiveAccount")

	accountName := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if accountName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountsClient := cognitiveservices.NewAccountsClient(subscriptionID)
	accountsClient.Authorizer = session.Authorizer

	account, err := accountsClient.Get(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if account.ID != nil {
		privateEndpointConnectionInfo := getPrivateEndpointConnectionInfo(&account)
		return &AccountInfo{account, privateEndpointConnectionInfo}, nil
	}

	return nil, nil
}

func listCognitiveAccountDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listCognitiveAccountDiagnosticSettings")
	id := *h.Item.(AccountInfo).Account.ID

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

func getPrivateEndpointConnectionInfo(account *cognitiveservices.Account) []PrivateEndpointConnectionInfo {
	privateEndpointConnectionInfo := []PrivateEndpointConnectionInfo{}
	if account.Properties.PrivateEndpointConnections != nil {
		for _, connection := range *account.Properties.PrivateEndpointConnections {
			privateEndpointConnectionPropertiesInfo := PrivateEndpointConnectionPropertiesInfo{}
			privateEndpointInfo := PrivateEndpointInfo{}
			if connection.Properties != nil {
				privateEndpointInfo.ID = connection.Properties.PrivateEndpoint.ID
				privateEndpointConnectionPropertiesInfo.PrivateEndpoint = privateEndpointInfo
				privateEndpointConnectionPropertiesInfo.PrivateLinkServiceConnectionState = connection.Properties.PrivateLinkServiceConnectionState
				privateEndpointConnectionPropertiesInfo.ProvisioningState = connection.Properties.ProvisioningState
				privateEndpointConnectionPropertiesInfo.GroupIds = connection.Properties.GroupIds
			}
			privateEndpointConnectionInfo = append(
				privateEndpointConnectionInfo,
				PrivateEndpointConnectionInfo {
					Properties: privateEndpointConnectionPropertiesInfo,
					SystemData: connection.SystemData,
					Location:   connection.Location,
					Etag:       connection.Etag,
					ID:         connection.ID,
					Name:       connection.Name,
					Type:       connection.Type,
				})
			}
		}
	
	return privateEndpointConnectionInfo
}
