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
				Transform:   transform.FromField("ApplicationGateway.Name"),
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGateway.ID"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the application gateway resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGateway.Type"),
			},
			{
				Name:        "enable_fips",
				Description: "Whether FIPS is enabled on the application gateway resource.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnableFips"),
				Default:     false,
			},
			{
				Name:        "enable_http2",
				Description: "Whether HTTP2 is enabled on the application gateway resource.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnableHTTP2"),
				Default:     false,
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGateway.Etag"),
			},
			{
				Name:        "force_firewall_policy_association",
				Description: "If true, associates a firewall policy with an application gateway regardless whether the policy differs from the WAF Config.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.ForceFirewallPolicyAssociation"),
				Default:     false,
			},
			{
				Name:        "operational_state",
				Description: "Operational state of the application gateway resource. Possible values include: 'Stopped', 'Starting', 'Running', 'Stopping'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.OperationalState"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the application gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ResourceGUID"),
			},
			{
				Name:        "authentication_certificates",
				Description: "Authentication certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AuthenticationCertificates"),
			},
			{
				Name:        "autoscale_configuration",
				Description: "Autoscale Configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AutoscaleConfiguration"),
			},
			{
				Name:        "backend_address_pools",
				Description: "Backend address pool of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.BackendAddressPools"),
			},
			{
				Name:        "backend_http_settings_collection",
				Description: "Backend http settings of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.BackendHTTPSettingsCollection"),
			},
			{
				Name:        "custom_error_configurations",
				Description: "Custom error configurations of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CustomErrorConfigurations"),
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
				Transform:   transform.FromField("Properties.FirewallPolicy"),
			},
			{
				Name:        "frontend_ip_configurations",
				Description: "Frontend IP addresses of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.FrontendIPConfigurations"),
			},
			{
				Name:        "frontend_ports",
				Description: "Frontend ports of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.FrontendPorts"),
			},
			{
				Name:        "gateway_ip_configurations",
				Description: "Subnets of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.GatewayIPConfigurations"),
			},
			{
				Name:        "http_listeners",
				Description: "Http listeners of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.HTTPListeners"),
			},
			{
				Name:        "identity",
				Description: "The identity of the application gateway, if configured.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGateway.Identity"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "Private endpoint connections on application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PrivateEndpointConnections"),
			},
			{
				Name:        "private_link_configurations",
				Description: "PrivateLink configurations on application gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PrivateLinkConfigurations"),
			},
			{
				Name:        "probes",
				Description: "Probes of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Probes"),
			},
			{
				Name:        "redirect_configurations",
				Description: "Redirect configurations of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.RedirectConfigurations"),
			},
			{
				Name:        "request_routing_rules",
				Description: "Request routing rules of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.RequestRoutingRules"),
			},
			{
				Name:        "rewrite_rule_sets",
				Description: "Rewrite rules for the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.RewriteRuleSets"),
			},
			{
				Name:        "sku",
				Description: "SKU of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Sku"),
			},
			{
				Name:        "ssl_certificates",
				Description: "SSL certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.SslCertificates"),
			},
			{
				Name:        "ssl_policy",
				Description: "SSL policy of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.SslPolicy"),
			},
			{
				Name:        "ssl_profiles",
				Description: "SSL profiles of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.SslProfiles"),
			},
			{
				Name:        "trusted_client_certificates",
				Description: "Trusted client certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.TrustedClientCertificates"),
			},
			{
				Name:        "trusted_root_certificates",
				Description: "Trusted root certificates of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.TrustedRootCertificates"),
			},
			{
				Name:        "url_path_maps",
				Description: "URL path map of the application gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.URLPathMaps"),
			},
			{
				Name:        "web_application_firewall_configuration",
				Description: "Web application firewall configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.WebApplicationFirewallConfiguration"),
			},
			{
				Name:        "zones",
				Description: "A list of availability zones denoting where the resource needs to come from.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGateway.Zones"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGateway.Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGateway.Tags"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationGateway.ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGateway.Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGateway.ID").Transform(extractResourceGroupFromID),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationGateway.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

type ApplicationGatewayInfo struct {
	ApplicationGateway network.ApplicationGateway
	Properties         ApplicationGatewayPropertiesInfo
}

type ApplicationGatewayPropertiesInfo struct {
	Sku 								interface{}
	SslPolicy 							interface{}
	OperationalState 					interface{}
	GatewayIPConfigurations 			interface{}
	AuthenticationCertificates 			interface{}
	TrustedRootCertificates 			interface{}
	TrustedClientCertificates 			interface{}
	SslCertificates 					interface{}
	FrontendIPConfigurations 			interface{}
	FrontendPorts 						interface{}
	Probes 								interface{}
	BackendAddressPools 				interface{}
	BackendHTTPSettingsCollection       interface{}
	HTTPListeners                       interface{}
	SslProfiles                         interface{}
	URLPathMaps                         interface{}
	RequestRoutingRules                 interface{}
	RewriteRuleSets                     interface{}
	RedirectConfigurations              interface{}
	WebApplicationFirewallConfiguration interface{}
	FirewallPolicy                      interface{}
	EnableHTTP2                         *bool
	EnableFips                          *bool
	AutoscaleConfiguration              interface{}
	PrivateLinkConfigurations           interface{}
	PrivateEndpointConnections          interface{}
	ResourceGUID                        *string
	ProvisioningState                   interface{}
	CustomErrorConfigurations           interface{}
	ForceFirewallPolicyAssociation      *bool
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
		applicationGatewayPropertiesInfo := getApplicationGatewayPropertiesInfo(&gateway)
		d.StreamListItem(ctx, ApplicationGatewayInfo{gateway, applicationGatewayPropertiesInfo})
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, gateway := range result.Values() {
			applicationGatewayPropertiesInfo := getApplicationGatewayPropertiesInfo(&gateway)
			d.StreamListItem(ctx, ApplicationGatewayInfo{gateway, applicationGatewayPropertiesInfo})
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
		applicationGatewayPropertiesInfo := getApplicationGatewayPropertiesInfo(&gateway)
		return &ApplicationGatewayInfo{gateway, applicationGatewayPropertiesInfo}, nil
	}

	return nil, nil
}

func listApplicationGatewayDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listApplicationGatewayDiagnosticSettings")
	id := *h.Item.(ApplicationGatewayInfo).ApplicationGateway.ID

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

func getApplicationGatewayPropertiesInfo(gateway *network.ApplicationGateway) ApplicationGatewayPropertiesInfo {
	applicationGatewayPropertiesInfo := ApplicationGatewayPropertiesInfo{}
	if gateway.ApplicationGatewayPropertiesFormat != nil {
		applicationGatewayPropertiesInfo.Sku =  gateway.ApplicationGatewayPropertiesFormat.Sku
		applicationGatewayPropertiesInfo.SslPolicy =  gateway.ApplicationGatewayPropertiesFormat.SslPolicy
		applicationGatewayPropertiesInfo.OperationalState =  gateway.ApplicationGatewayPropertiesFormat.OperationalState
		applicationGatewayPropertiesInfo.FrontendIPConfigurations =  gateway.ApplicationGatewayPropertiesFormat.FrontendIPConfigurations
		applicationGatewayPropertiesInfo.FrontendPorts =  gateway.ApplicationGatewayPropertiesFormat.FrontendPorts
		applicationGatewayPropertiesInfo.Probes =  gateway.ApplicationGatewayPropertiesFormat.Probes
		applicationGatewayPropertiesInfo.BackendAddressPools =  gateway.ApplicationGatewayPropertiesFormat.BackendAddressPools
		applicationGatewayPropertiesInfo.BackendHTTPSettingsCollection =  gateway.ApplicationGatewayPropertiesFormat.BackendHTTPSettingsCollection
		applicationGatewayPropertiesInfo.HTTPListeners =  gateway.ApplicationGatewayPropertiesFormat.HTTPListeners
		applicationGatewayPropertiesInfo.SslProfiles =  gateway.ApplicationGatewayPropertiesFormat.SslProfiles
		applicationGatewayPropertiesInfo.URLPathMaps =  gateway.ApplicationGatewayPropertiesFormat.URLPathMaps
		applicationGatewayPropertiesInfo.RequestRoutingRules =  gateway.ApplicationGatewayPropertiesFormat.RequestRoutingRules
		applicationGatewayPropertiesInfo.RewriteRuleSets =  gateway.ApplicationGatewayPropertiesFormat.RewriteRuleSets
		applicationGatewayPropertiesInfo.RedirectConfigurations =  gateway.ApplicationGatewayPropertiesFormat.RedirectConfigurations
		applicationGatewayPropertiesInfo.WebApplicationFirewallConfiguration =  gateway.ApplicationGatewayPropertiesFormat.WebApplicationFirewallConfiguration
		applicationGatewayPropertiesInfo.FirewallPolicy =  gateway.ApplicationGatewayPropertiesFormat.FirewallPolicy
		applicationGatewayPropertiesInfo.AutoscaleConfiguration =  gateway.ApplicationGatewayPropertiesFormat.AutoscaleConfiguration
		applicationGatewayPropertiesInfo.PrivateLinkConfigurations =  gateway.ApplicationGatewayPropertiesFormat.PrivateLinkConfigurations
		applicationGatewayPropertiesInfo.PrivateEndpointConnections =  gateway.ApplicationGatewayPropertiesFormat.PrivateEndpointConnections
		applicationGatewayPropertiesInfo.ResourceGUID =  gateway.ApplicationGatewayPropertiesFormat.ResourceGUID
		applicationGatewayPropertiesInfo.ProvisioningState =  gateway.ApplicationGatewayPropertiesFormat.ProvisioningState
		applicationGatewayPropertiesInfo.CustomErrorConfigurations =  gateway.ApplicationGatewayPropertiesFormat.CustomErrorConfigurations
		applicationGatewayPropertiesInfo.ForceFirewallPolicyAssociation =  gateway.ApplicationGatewayPropertiesFormat.ForceFirewallPolicyAssociation
		applicationGatewayPropertiesInfo.GatewayIPConfigurations = gateway.ApplicationGatewayPropertiesFormat.GatewayIPConfigurations
		applicationGatewayPropertiesInfo.AuthenticationCertificates =  gateway.ApplicationGatewayPropertiesFormat.AuthenticationCertificates
		applicationGatewayPropertiesInfo.TrustedRootCertificates =  gateway.ApplicationGatewayPropertiesFormat.TrustedRootCertificates
		applicationGatewayPropertiesInfo.TrustedClientCertificates =  gateway.ApplicationGatewayPropertiesFormat.TrustedClientCertificates
		applicationGatewayPropertiesInfo.SslCertificates =  gateway.ApplicationGatewayPropertiesFormat.SslCertificates
		applicationGatewayPropertiesInfo.EnableHTTP2 =  gateway.ApplicationGatewayPropertiesFormat.EnableHTTP2
		applicationGatewayPropertiesInfo.EnableFips =  gateway.ApplicationGatewayPropertiesFormat.EnableFips
	}

	return applicationGatewayPropertiesInfo
}
