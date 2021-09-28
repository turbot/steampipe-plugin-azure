package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/healthcareapis/mgmt/healthcareapis"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureHealthcareService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_healthcare_service",
		Description: "Azure Healthcare Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getHealthcareService,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listHealthcareServices,
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
				Name:        "etag",
				Description: "An etag associated with the resource, used for optimistic concurrency when editing it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the healthcare service resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "allow_credentials",
				Description: "If credentials are allowed via CORS.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.CorsConfiguration.AllowCredentials"),
			},
			{
				Name:        "audience",
				Description: "The audience url for the service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AuthenticationConfiguration.Audience").Transform(transform.ToString),
			},
			{
				Name:        "authority",
				Description: "The authority url for the service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AuthenticationConfiguration.Authority").Transform(transform.ToString),
			},
			{
				Name:        "cosmos_db_configuration",
				Description: "The settings for the Cosmos DB database backing the service.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.CosmosDbConfiguration.OfferThroughput"),
			},
			{
				Name:        "kind",
				Description: "The kind of the service. Possible values include: 'Fhir', 'FhirStu3', 'FhirR4'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The resource location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_age",
				Description: "The max age to be allowed via CORS.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.CorsConfiguration.MaxAge"),
			},
			{
				Name:        "smart_proxy_enabled",
				Description: "If the SMART on FHIR proxy is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.AuthenticationConfiguration.SmartProxyEnabled"),
			},
			{
				Name:        "access_policies",
				Description: "The access policies of the healthcare service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AccessPolicies"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the healthcare serive.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHealthcareServiceDignosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "headers",
				Description: "The headers to be allowed via CORS.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CorsConfiguration.Origins"),
			},
			{
				Name:        "methods",
				Description: "The methods to be allowed via CORS.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CorsConfiguration.Methods"),
			},
			{
				Name:        "origins",
				Description: "The origins to be allowed via CORS.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CorsConfiguration.Origins"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "List of private endpoint connections for healthcare service.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHealthcarePrivateEndpointConnections,
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

func listHealthcareServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	healthcareClient := healthcareapis.NewServicesClient(subscriptionID)
	healthcareClient.Authorizer = session.Authorizer
	result, err := healthcareClient.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listHealthcareServices", "list", err)
		return nil, err
	}

	for _, service := range result.Values() {
		d.StreamListItem(ctx, service)
	}

	for result.NotDone() {
		err := result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listHealthcareServices", "paging", err)
			return nil, err
		}

		for _, service := range result.Values() {
			d.StreamListItem(ctx, service)
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getHealthcareService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHealthcareService")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Empty check for param
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	serviceClient := healthcareapis.NewServicesClient(subscriptionID)
	serviceClient.Authorizer = session.Authorizer

	op, err := serviceClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getHealthcareService", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

// If we return the API response directly, the output will not provide the properties of PrivateEndpointConnections
func getHealthcarePrivateEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHealthcarePrivateEndpointConnections")

	serviceDetails := h.Item.(healthcareapis.ServicesDescription)

	// Empty check
	if serviceDetails.ID == nil || serviceDetails.Name == nil {
		return nil, nil
	}

	resourceGroup := strings.Split(*serviceDetails.ID, "/")[4]
	resourceName := serviceDetails.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	serviceClient := healthcareapis.NewPrivateEndpointConnectionsClient(subscriptionID)
	serviceClient.Authorizer = session.Authorizer

	// SDK does not support pagination yet
	op, err := serviceClient.ListByService(ctx, resourceGroup, *resourceName)

	if err != nil {
		plugin.Logger(ctx).Error("getHealthcarePrivateEndpointConnections", "list", err)
		return nil, err
	}

	var privateEndpoints []map[string]interface{}
	
	for _, conn := range *op.Value {
		privateEndpoint := make(map[string]interface{})
		if conn.ID != nil {
			privateEndpoint["PrivateEndpointConnectionId"] = conn.ID
		}
		if conn.Name != nil {
			privateEndpoint["PrivateEndpointConnectionName"] = conn.Name
		}
		if conn.Type != nil {
			privateEndpoint["PrivateEndpointConnectionType"] = conn.Type
		}
		if conn.PrivateEndpointConnectionProperties != nil {
			if conn.PrivateEndpointConnectionProperties.PrivateEndpoint != nil {
				if conn.PrivateEndpointConnectionProperties.PrivateEndpoint.ID != nil {
					privateEndpoint["PrivateEndpointId"] = conn.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
				}
			}
			if conn.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState != nil {
				if conn.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired != nil {
					privateEndpoint["PrivateLinkServiceConnectionStateActionsRequired"] = conn.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired
				}
				if conn.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status != "" {
					privateEndpoint["PrivateLinkServiceConnectionStateStatus"] = conn.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status
				}
				if conn.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description != nil {
					privateEndpoint["PrivateLinkServiceConnectionStateDescription"] = conn.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description
				}
			}
			if conn.PrivateEndpointConnectionProperties.ProvisioningState != "" {
				privateEndpoint["ProvisioningState"] = conn.PrivateEndpointConnectionProperties.ProvisioningState
			}
		}
		privateEndpoints = append(privateEndpoints, privateEndpoint)
	}

	return privateEndpoints, nil
}

func getHealthcareServiceDignosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHealthcareServiceDignosisSettings")

	serviceDetails := h.Item.(healthcareapis.ServicesDescription)

	// Empty check
	if serviceDetails.ID == nil {
		return nil, nil
	}

	resourceId := serviceDetails.ID

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	dignosticSettingClient := insights.NewDiagnosticSettingsClient(subscriptionID)
	dignosticSettingClient.Authorizer = session.Authorizer

	op, err := dignosticSettingClient.List(ctx, *resourceId)
	if err != nil {
		plugin.Logger(ctx).Error("getHealthcareServiceDignosisSettings", "list", err)
		return nil, err
	}

	// If we return the API response directly, the output will not provide all
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
