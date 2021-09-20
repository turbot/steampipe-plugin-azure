package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureApplicationGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_application_gateway",
		Description: "Azure Application Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getApplicationGateway,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listApplicationGateways,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the application gateway resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "enable_fips",
				Description: "Whether FIPS is enabled on the application gateway resource.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.EnableFips"),
				Default:     false,
			},
			{
				Name:        "enable_http2",
				Description: "Whether HTTP2 is enabled on the application gateway resource.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.EnableHTTP2"),
				Default:     false,
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Etag"),
			},
			{
				Name:        "force_firewall_policy_association",
				Description: "If true, associates a firewall policy with an application gateway regardless whether the policy differs from the WAF Config.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.ForceFirewallPolicyAssociation"),
				Default:     false,
			},
			{
				Name:        "operational_state",
				Description: "Operational state of the application gateway resource. Possible values include: 'Stopped', 'Starting', 'Running', 'Stopping'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.OperationalState"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the application gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "authentication_certificates",
				Description: "Authentication certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractAuthenticationCertificates),
			},
			{
				Name:        "autoscale_configuration",
				Description: "Autoscale Configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.AutoscaleConfiguration"),
			},
			{
				Name:        "backend_address_pools",
				Description: "Backend address pool of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractBackendAddressPools),
			},
			{
				Name:        "backend_http_settings_collection",
				Description: "Backend http settings of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractBackendHTTPSettingsCollection),
			},
			{
				Name:        "custom_error_configurations",
				Description: "Custom error configurations of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.CustomErrorConfigurations"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the application gateway.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listApplicationGatewayDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "firewall_policy",
				Description: "Reference to the FirewallPolicy resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.FirewallPolicy"),
			},
			{
				Name:        "frontend_ip_configurations",
				Description: "Frontend IP addresses of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractFrontendIPConfigurations),
			},
			{
				Name:        "frontend_ports",
				Description: "Frontend ports of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractFrontendPorts),
			},
			{
				Name:        "gateway_ip_configurations",
				Description: "Subnets of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayIPConfigurations),
			},
			{
				Name:        "http_listeners",
				Description: "Http listeners of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractHTTPListeners),
			},
			{
				Name:        "identity",
				Description: "The identity of the application gateway, if configured.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractIdentity),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "Private endpoint connections on application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractPrivateEndpointConnections),
			},
			{
				Name:        "private_link_configurations",
				Description: "PrivateLink configurations on application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractPrivateLinkConfigurations),
			},
			{
				Name:        "probes",
				Description: "Probes of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractProbes),
			},
			{
				Name:        "redirect_configurations",
				Description: "Redirect configurations of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.RedirectConfigurations"),
			},
			{
				Name:        "request_routing_rules",
				Description: "Request routing rules of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractRequestRoutingRules),
			},
			{
				Name:        "rewrite_rule_sets",
				Description: "Rewrite rules for the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractRewriteRuleSets),
			},
			{
				Name:        "sku",
				Description: "SKU of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.Sku"),
			},
			{
				Name:        "ssl_certificates",
				Description: "SSL certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSslCertificates),
			},
			{
				Name:        "ssl_policy",
				Description: "SSL policy of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.SslPolicy"),
			},
			{
				Name:        "ssl_profiles",
				Description: "SSL profiles of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSslProfiles),
			},
			{
				Name:        "trusted_client_certificates",
				Description: "Trusted client certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractTrustedClientCertificates),
			},
			{
				Name:        "trusted_root_certificates",
				Description: "Trusted root certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractTrustedRootCertificates),
			},
			{
				Name:        "url_path_maps",
				Description: "URL path map of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractURLPathMaps),
			},
			{
				Name:        "web_application_firewall_configuration",
				Description: "Web application firewall configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.WebApplicationFirewallConfiguration"),
			},
			{
				Name:        "zones",
				Description: "A list of availability zones denoting where the resource needs to come from.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Zones"),
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

type Identity struct {
	PrincipalID             *string
	TenantID                *string
	Type                    interface{}
	UserAssignedIdentities []UserAssignedIdentitiesValue
}

type UserAssignedIdentitiesValue struct {
	PrincipalID *string
	ClientID    *string
}

//// LIST FUNCTION

func listApplicationGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := network.NewApplicationGatewaysClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListAll(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listApplicationGateways", "list", err)
		return nil, err
	}

	for _, gateway := range result.Values() {
		d.StreamListItem(ctx, gateway)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listApplicationGateways", "list_paging", err)
			return nil, err
		}
		for _, gateway := range result.Values() {
			d.StreamListItem(ctx, gateway)
		}
	}
	
	return nil, err
}

//// HYDRATE FUNCTIONS

func getApplicationGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getApplicationGateway")

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

	client := network.NewApplicationGatewaysClient(subscriptionID)
	client.Authorizer = session.Authorizer

	gateway, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getApplicationGateway", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if gateway.ID != nil {
		return gateway, nil
	}

	return nil, nil
}

func listApplicationGatewayDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listApplicationGatewayDiagnosticSettings")
	id := *h.Item.(network.ApplicationGateway).ID

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

// If we return the API response directly, the output will not provide all the properties of Identity
func extractIdentity(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var identity Identity
	if gateway.Identity != nil {
		identity.PrincipalID = gateway.Identity.PrincipalID
		identity.TenantID = gateway.Identity.TenantID
		identity.Type = gateway.Identity.Type
		for _, userAssignedIdentity := range gateway.Identity.UserAssignedIdentities {
			var userAssignedIdentitiesValue UserAssignedIdentitiesValue
			if userAssignedIdentity.ClientID != nil {
				userAssignedIdentitiesValue.ClientID = userAssignedIdentity.ClientID
			}
			if userAssignedIdentity.PrincipalID != nil {
				userAssignedIdentitiesValue.PrincipalID = userAssignedIdentity.PrincipalID
			}

			identity.UserAssignedIdentities = append(identity.UserAssignedIdentities, userAssignedIdentitiesValue)
		}
	}
	return identity, nil
}

// If we return the API response directly, the output will not provide all the properties of GatewayIPConfigurations
func extractGatewayIPConfigurations(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.GatewayIPConfigurations != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.GatewayIPConfigurations {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayIPConfigurationPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayIPConfigurationPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayIPConfigurationPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of AuthenticationCertificates
func extractAuthenticationCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.AuthenticationCertificates != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.AuthenticationCertificates {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayAuthenticationCertificatePropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayAuthenticationCertificatePropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayAuthenticationCertificatePropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of TrustedRootCertificates
func extractTrustedRootCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.TrustedRootCertificates != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.TrustedRootCertificates {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayTrustedRootCertificatePropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayTrustedRootCertificatePropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayTrustedRootCertificatePropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of TrustedClientCertificates
func extractTrustedClientCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.TrustedClientCertificates != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.TrustedClientCertificates {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayTrustedClientCertificatePropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayTrustedClientCertificatePropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayTrustedClientCertificatePropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of SslCertificates
func extractSslCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.SslCertificates != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.SslCertificates {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewaySslCertificatePropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewaySslCertificatePropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewaySslCertificatePropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of FrontendIPConfigurations
func extractFrontendIPConfigurations(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.FrontendIPConfigurations != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.FrontendIPConfigurations {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayFrontendIPConfigurationPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayFrontendIPConfigurationPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayFrontendIPConfigurationPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of FrontendPorts
func extractFrontendPorts(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.FrontendPorts != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.FrontendPorts {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayFrontendPortPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayFrontendPortPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayFrontendPortPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of Probes
func extractProbes(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.Probes != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.Probes {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayProbePropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayProbePropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayProbePropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of BackendAddressPools
func extractBackendAddressPools(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.BackendAddressPools != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.BackendAddressPools {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayBackendAddressPoolPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayBackendAddressPoolPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayBackendAddressPoolPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of BackendHTTPSettingsCollection
func extractBackendHTTPSettingsCollection(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.BackendHTTPSettingsCollection != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.BackendHTTPSettingsCollection {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayBackendHTTPSettingsPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayBackendHTTPSettingsPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of HTTPListeners
func extractHTTPListeners(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.HTTPListeners != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.HTTPListeners {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayHTTPListenerPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayHTTPListenerPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayHTTPListenerPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of SslProfiles
func extractSslProfiles(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.SslProfiles != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.SslProfiles {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewaySslProfilePropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewaySslProfilePropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewaySslProfilePropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of URLPathMaps
func extractURLPathMaps(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.URLPathMaps != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.URLPathMaps {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayURLPathMapPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayURLPathMapPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayURLPathMapPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of RequestRoutingRules
func extractRequestRoutingRules(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.RequestRoutingRules != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.RequestRoutingRules {
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
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayRequestRoutingRulePropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayRequestRoutingRulePropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayRequestRoutingRulePropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of RewriteRuleSets
func extractRewriteRuleSets(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.RewriteRuleSets != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.RewriteRuleSets {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayRewriteRuleSetPropertiesFormat != nil {
				objectMap["properties"] = i.ApplicationGatewayRewriteRuleSetPropertiesFormat
				objectMap["provisioning_state"] = i.ApplicationGatewayRewriteRuleSetPropertiesFormat.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of PrivateLinkConfigurations
func extractPrivateLinkConfigurations(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.PrivateLinkConfigurations != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.PrivateLinkConfigurations {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayPrivateLinkConfigurationProperties != nil {
				objectMap["properties"] = i.ApplicationGatewayPrivateLinkConfigurationProperties
				objectMap["provisioning_state"] = i.ApplicationGatewayPrivateLinkConfigurationProperties.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}

// If we return the API response directly, the output will not provide all the properties of PrivateEndpointConnections
func extractPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.ApplicationGateway)
	var properties []map[string]interface{}

	if gateway.ApplicationGatewayPropertiesFormat.PrivateEndpointConnections != nil {
		for _, i := range *gateway.ApplicationGatewayPropertiesFormat.PrivateEndpointConnections {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Etag != nil {
				objectMap["type"] = i.Etag
			}
			if i.ApplicationGatewayPrivateEndpointConnectionProperties != nil {
				objectMap["properties"] = i.ApplicationGatewayPrivateEndpointConnectionProperties
				objectMap["provisioning_state"] = i.ApplicationGatewayPrivateEndpointConnectionProperties.ProvisioningState
			}
			properties = append(properties, objectMap)
		}
	}
	
	return properties, nil
}
