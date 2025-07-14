package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureAppServiceFunctionApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_function_app",
		Description: "Azure App Service Function App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAppServiceFunctionApp,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "sites/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAppServiceFunctionApps,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "sites/read",
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAppServiceFunctionAppSiteConfiguration,
				Tags: map[string]string{
					"service": "Microsoft.Web",
					"action":  "sites/config/read",
				},
			},
			{
				Func: getAppServiceFunctionAppSiteAuthSetting,
				Tags: map[string]string{
					"service": "Microsoft.Web",
					"action":  "sites/config/authsettings/read",
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the app service function app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an app service function app uniquely.",
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
				Description: "The resource type of the app service function app.",
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
				Hydrate:     getAppServiceFunctionAppSiteAuthSetting,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "configuration",
				Description: "Describes the configuration of an app.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppServiceFunctionAppSiteConfiguration,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "site_config",
				Description: "A map of all configuration for the app",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.SiteConfig"),
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

func listAppServiceFunctionApps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &webClient, d.Connection)

	result, err := webClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, functionApp := range result.Values() {
		// Filtering out all the web apps
		if strings.Contains(string(*functionApp.Kind), "functionapp") {
			d.StreamListItem(ctx, functionApp)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	for result.NotDone() {
		err := result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, functionApp := range result.Values() {
			// Filtering out all the web apps
			if strings.Contains(string(*functionApp.Kind), "functionapp") {
				d.StreamListItem(ctx, functionApp)
				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAppServiceFunctionApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceFunctionApp")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

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

	// Apply Retry rule
	ApplyRetryRules(ctx, &webClient, d.Connection)

	op, err := webClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil && strings.Contains(string(*op.Kind), "functionapp") {
		return op, nil
	}

	return nil, nil
}

func getAppServiceFunctionAppSiteConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceFunctionAppSiteConfiguration")

	data := h.Item.(web.Site)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &webClient, d.Connection)

	op, err := webClient.GetConfiguration(ctx, *data.SiteProperties.ResourceGroup, *data.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getAppServiceFunctionAppSiteAuthSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceFunctionAppSiteAuthSetting")

	data := h.Item.(web.Site)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &webClient, d.Connection)

	op, err := webClient.GetAuthSettings(ctx, *data.SiteProperties.ResourceGroup, *data.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
