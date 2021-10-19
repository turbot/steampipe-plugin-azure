package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/appconfiguration/mgmt/2020-06-01/appconfiguration"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAppConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_configuration",
		Description: "Azure App Configuration",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAppConfiguration,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAppConfigurations,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the configuration store. Possible values include: 'Creating', 'Updating', 'Deleting', 'Succeeded', 'Failed', 'Canceled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigurationStoreProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The creation date of configuration store.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ConfigurationStoreProperties.CreationDate").Transform(convertDateToTime),
			},
			{
				Name:        "endpoint",
				Description: "The DNS endpoint where the configuration store API will be available.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigurationStoreProperties.Endpoint"),
			},
			{
				Name:        "public_network_access",
				Description: "Control permission for data plane traffic coming from public networks while private endpoint is enabled. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getPublicNetworkAccess,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "sku_name",
				Description: "The SKU name of the configuration store.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the configuration store.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAppConfigurationDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "The encryption settings of the configuration store.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConfigurationStoreProperties.Encryption"),
			},
			{
				Name:        "identity",
				Description: "The managed identity information, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The list of private endpoint connections that are set up for this resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractAppConfigurationPrivateEndpointConnections),
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
				Transform:   transform.FromField("Tags"),
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

func listAppConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := appconfiguration.NewConfigurationStoresClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx, "")
	if err != nil {
		plugin.Logger(ctx).Error("listAppConfigurations", "list", err)
		return nil, err
	}

	for _, config := range result.Values() {
		d.StreamListItem(ctx, config)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listAppConfigurations", "list_paging", err)
			return nil, err
		}
		for _, config := range result.Values() {
			d.StreamListItem(ctx, config)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAppConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppConfiguration")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := appconfiguration.NewConfigurationStoresClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	config, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getAppConfiguration", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if config.ID != nil {
		return config, nil
	}

	return nil, nil
}

func listAppConfigurationDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAppConfigurationDiagnosticSettings")
	id := *h.Item.(appconfiguration.ConfigurationStore).ID

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
		plugin.Logger(ctx).Error("listAppConfigurationDiagnosticSettings", "list", err)
		return nil, err
	}

	// If we return the API response directly, the output does not provide all
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

func getPublicNetworkAccess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	configurationStore := h.Item.(appconfiguration.ConfigurationStore)

	// In case of automatic is selected at the time of store creation, PublicNetworkAccess value will be nil in API response.
	// With a private endpoint, public network access will be automatically disabled.
	// If there is no private endpoint present, public network access is automatically enabled.
	if len(configurationStore.PublicNetworkAccess) == 0 && len(*configurationStore.PrivateEndpointConnections) == 0 {
		return "Enabled", nil
	} else if len(configurationStore.PublicNetworkAccess) == 0 && len(*configurationStore.PrivateEndpointConnections) != 0 {
		return "Disabled", nil
	}

	return configurationStore.PublicNetworkAccess, nil
}

//// TRANSFORM FUNCTION

// If we return the API response directly, the output will not provide all the properties of PrivateEndpointConnections
func extractAppConfigurationPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	server := d.HydrateItem.(appconfiguration.ConfigurationStore)
	var properties []map[string]interface{}

	if server.ConfigurationStoreProperties.PrivateEndpointConnections != nil {
		for _, i := range *server.ConfigurationStoreProperties.PrivateEndpointConnections {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.ID != nil {
				objectMap["name"] = i.Name
			}
			if i.ID != nil {
				objectMap["type"] = i.Type
			}
			if i.PrivateEndpointConnectionProperties != nil {
				if i.PrivateEndpointConnectionProperties.PrivateEndpoint != nil {
					objectMap["privateEndpointPropertyId"] = i.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
				}
				if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState != nil {
					if len(i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired) > 0 {
						objectMap["privateLinkServiceConnectionStateActionsRequired"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired
					}
					if len(i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status) > 0 {
						objectMap["privateLinkServiceConnectionStateStatus"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status
					}
					if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description != nil {
						objectMap["privateLinkServiceConnectionStateDescription"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description
					}
				}
				if len(i.PrivateEndpointConnectionProperties.ProvisioningState) > 0 {
					objectMap["provisioningState"] = i.PrivateEndpointConnectionProperties.ProvisioningState
				}
			}
			properties = append(properties, objectMap)
		}
	}

	return properties, nil
}
