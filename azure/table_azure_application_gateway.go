package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureApplicationGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_application_gateway",
		Description: "Azure Application Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getApplicationGateway,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listApplicationGateways,
		},
		Columns: azureColumns([]*plugin.Column{
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
				Description: "The provisioning state of the application gateway. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
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
				Description: "Whether FIPS is enabled on the application gateway.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.EnableFips"),
				Default:     false,
			},
			{
				Name:        "enable_http2",
				Description: "Whether HTTP2 is enabled on the application gateway.",
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
				Description: "If true, associates a firewall policy with an application gateway regardless whether the policy differs from the WAF configuration.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.ForceFirewallPolicyAssociation"),
				Default:     false,
			},
			{
				Name:        "operational_state",
				Description: "Operational state of the application gateway. Possible values include: 'Stopped', 'Starting', 'Running', 'Stopping'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.OperationalState"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the application gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "authentication_certificates",
				Description: "Authentication certificates of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayAuthenticationCertificates),
			},
			{
				Name:        "autoscale_configuration",
				Description: "Autoscale Configuration of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.AutoscaleConfiguration"),
			},
			{
				Name:        "backend_address_pools",
				Description: "Backend address pool of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayBackendAddressPools),
			},
			{
				Name:        "backend_http_settings_collection",
				Description: "Backend http settings of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayBackendHTTPSettingsCollection),
			},
			{
				Name:        "custom_error_configurations",
				Description: "Custom error configurations of the application gateway.",
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
			// This column value will be populated if the firewall policy is associated and the firewall configuration is not disabled for the Application Gateway.
			{
				Name:        "firewall_policy",
				Description: "Reference to the FirewallPolicy resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.FirewallPolicy"),
			},
			{
				Name:        "frontend_ip_configurations",
				Description: "Frontend IP addresses of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayFrontendIPConfigurations),
			},
			{
				Name:        "frontend_ports",
				Description: "Frontend ports of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayFrontendPorts),
			},
			{
				Name:        "gateway_ip_configurations",
				Description: "Subnets of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayIPConfigurations),
			},
			{
				Name:        "http_listeners",
				Description: "Http listeners of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayHTTPListeners),
			},
			{
				Name:        "identity",
				Description: "The identity of the application gateway, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_endpoint_connections",
				Description: "Private endpoint connections on application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayPrivateEndpointConnections),
			},
			{
				Name:        "private_link_configurations",
				Description: "PrivateLink configurations on application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayPrivateLinkConfigurations),
			},
			{
				Name:        "probes",
				Description: "Probes of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayProbes),
			},
			{
				Name:        "redirect_configurations",
				Description: "Redirect configurations of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.RedirectConfigurations"),
			},
			{
				Name:        "request_routing_rules",
				Description: "Request routing rules of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayRequestRoutingRules),
			},
			{
				Name:        "rewrite_rule_sets",
				Description: "Rewrite rules for the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayRewriteRuleSets),
			},
			{
				Name:        "sku",
				Description: "SKU of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.Sku"),
			},
			{
				Name:        "ssl_certificates",
				Description: "SSL certificates of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewaySslCertificates),
			},
			{
				Name:        "ssl_policy",
				Description: "SSL policy of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGatewayPropertiesFormat.SslPolicy"),
			},
			{
				Name:        "ssl_profiles",
				Description: "SSL profiles of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewaySslProfiles),
			},
			{
				Name:        "trusted_client_certificates",
				Description: "Trusted client certificates of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayTrustedClientCertificates),
			},
			{
				Name:        "trusted_root_certificates",
				Description: "Trusted root certificates of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayTrustedRootCertificates),
			},
			{
				Name:        "url_path_maps",
				Description: "URL path map of the application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractGatewayURLPathMaps),
			},
			// This column value will be populated once the background configuration for the Application Gateway is complete. And if the tier 'WAF V2' is selected under the Settings > Configuration.
			{
				Name:        "web_application_firewall_configuration",
				Description: "Web application firewall configuration of the application gateway.",
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
		}),
	}
}

//// LIST FUNCTION

func listApplicationGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := network.NewApplicationGatewaysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := network.NewApplicationGatewaysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output does not provide
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

//// TRANSFORM FUNCTIONS

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
func extractGatewayAuthenticationCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayTrustedRootCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayTrustedClientCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewaySslCertificates(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayFrontendIPConfigurations(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayFrontendPorts(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayProbes(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayBackendAddressPools(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayBackendHTTPSettingsCollection(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayHTTPListeners(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewaySslProfiles(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayURLPathMaps(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayRequestRoutingRules(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayRewriteRuleSets(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayPrivateLinkConfigurations(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func extractGatewayPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
