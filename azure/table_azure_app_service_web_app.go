package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureAppServiceWebApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_web_app",
		Description: "Azure App Service Web App",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAppServiceWebApp,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAppServiceWebApps,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
			{
				Func:    getAppServiceWebAppVnetConnection,
				Depends: []plugin.HydrateFunc{getAppServiceWebAppSiteConfiguration},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the app service web app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an app service web app uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "kind",
				Description: "Contains the kind of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Current state of the app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.State"),
			},
			{
				Name:        "type",
				Description: "The resource type of the app service web app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_affinity_enabled",
				Description: "Specify whether client affinity is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.ClientAffinityEnabled"),
			},
			{
				Name:        "client_cert_enabled",
				Description: "Specify whether client certificate authentication is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.ClientCertEnabled"),
			},
			{
				Name:        "default_site_hostname",
				Description: "Default hostname of the app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.DefaultHostName"),
			},
			{
				Name:        "enabled",
				Description: "Specify whether the app is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.Enabled"),
			},
			{
				Name:        "host_name_disabled",
				Description: "Specify whether the public hostnames of the app is disabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.HostNamesDisabled"),
			},
			{
				Name:        "https_only",
				Description: "Specify whether configuring a web site to accept only https requests.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.HTTPSOnly"),
			},
			{
				Name:        "identity",
				Description: "Managed service identity for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(webAppIdentity),
			},
			{
				Name:        "outbound_ip_addresses",
				Description: "List of IP addresses that the app uses for outbound connections (e.g. database access).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.OutboundIPAddresses"),
			},
			{
				Name:        "possible_outbound_ip_addresses",
				Description: "List of possible IP addresses that the app uses for outbound connections (e.g. database access).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.PossibleOutboundIPAddresses"),
			},
			{
				Name:        "reserved",
				Description: "Specify whether the app is reserved.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.Reserved"),
			},
			{
				Name:        "host_names",
				Description: "A list of hostnames associated with the app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.HostNames"),
			},
			{
				Name:        "auth_settings",
				Description: "Describes the Authentication/Authorization settings of an app.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppServiceWebAppSiteAuthSetting,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "configuration",
				Description: "Describes the configuration of an app.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppServiceWebAppSiteConfiguration,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "site_config",
				Description: "A map of all configuration for the app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.SiteConfig"),
			},
			{
				Name:        "vnet_connection",
				Description: "Describes the virtual network connection for the app.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppServiceWebAppVnetConnection,
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
				Transform:   transform.FromField("SiteProperties.ResourceGroup").Transform(toLower),
			},
		}),
	}
}

//// LIST FUNCTION

func listAppServiceWebApps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	result, err := webClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, webApp := range result.Values() {
		// Filtering out all the function apps
		if string(*webApp.Kind) != "functionapp" {
			d.StreamListItem(ctx, webApp)
		}
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

		for _, webApp := range result.Values() {
			// Filtering out all the function apps
			if string(*webApp.Kind) != "functionapp" {
				d.StreamListItem(ctx, webApp)
			}
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

func getAppServiceWebApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceWebApp")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = web.AppsClient#Get: Invalid input: autorest/validation: validation failed: parameter=resourceGroupName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	op, err := webClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil && string(*op.Kind) != "functionapp" {
		return op, nil
	}

	return nil, nil
}

func getAppServiceWebAppSiteConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceWebAppSiteConfiguration")

	data := h.Item.(web.Site)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	op, err := webClient.GetConfiguration(ctx, *data.SiteProperties.ResourceGroup, *data.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getAppServiceWebAppSiteAuthSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceWebAppSiteAuthSetting")

	data := h.Item.(web.Site)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	op, err := webClient.GetAuthSettings(ctx, *data.SiteProperties.ResourceGroup, *data.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getAppServiceWebAppVnetConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceWebAppVnetConnection")

	data := h.Item.(web.Site)
	vnet := h.HydrateResults["getAppServiceWebAppSiteConfiguration"].(web.SiteConfigResource)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	// Return nil, if no virtual network is configured
	if *vnet.SiteConfig.VnetName == "" {
		return nil, nil
	}

	op, err := webClient.GetVnetConnection(ctx, *data.SiteProperties.ResourceGroup, *data.Name, *vnet.SiteConfig.VnetName)
	if err != nil {
		return nil, err
	}

	if op.VnetInfoProperties != nil {
		appVnetConnection := make(map[string]interface{})
		if op.Name != nil {
			appVnetConnection["name"] = op.Name
		}
		if op.ID != nil {
			appVnetConnection["id"] = op.ID
		}
		if op.Type != nil {
			appVnetConnection["type"] = op.Type
		}
		appVnetConnection["properties"] = op.VnetInfoProperties

		return appVnetConnection, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func webAppIdentity(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(web.Site)
	objectMap := make(map[string]interface{})
	if data.Identity != nil {
		if types.SafeString(data.Identity.Type) != "" {
			objectMap["Type"] = data.Identity.Type
		}
		if data.Identity.TenantID != nil {
			objectMap["TenantID"] = data.Identity.TenantID
		}
		if data.Identity.PrincipalID != nil {
			objectMap["PrincipalID"] = data.Identity.PrincipalID
		}
		if data.Identity.UserAssignedIdentities != nil {
			objectMap["UserAssignedIdentities"] = data.Identity.UserAssignedIdentities
		}
	}
	return objectMap, nil
}
