package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureAppServiceWebAppSlot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_web_app_slot",
		Description: "Azure App Service Web App Slot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "app_name", "resource_group"}),
			Hydrate:    getAppServiceWebAppSlot,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAppServiceWebApps,
			Hydrate:       listAppServiceWebAppSlots,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "app_name",
					Require: plugin.Optional,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Resource Name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "app_name",
				Description: "The name of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Resource ID of the app slot.",
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
				Description: "Resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_time_utc",
				Description: "Last time the app was modified, in UTC.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SiteProperties.LastModifiedTimeUtc.TIme"),
			},
			{
				Name:        "repository_site_name",
				Description: "Name of the repository site.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.RepositorySiteName"),
			},
			{
				Name:        "usage_state",
				Description: "State indicating whether the app has exceeded its quota usage. Read-only. Possible values include: 'UsageStateNormal', 'UsageStateExceeded'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.UsageState"),
			},
			{
				Name:        "enabled",
				Description: "Indicates wheather the app is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.Enabled"),
			},
			{
				Name:        "availability_state",
				Description: "Management information availability state for the app. Possible values include: 'Normal', 'Limited', 'DisasterRecoveryMode'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.AvailabilityState"),
			},
			{
				Name:        "server_farm_id",
				Description: "Resource ID of the associated App Service plan, formatted as: '/subscriptions/{subscriptionID}/resourceGroups/{groupName}/providers/Microsoft.Web/serverfarms/{appServicePlanName}'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.ServerFarmID"),
			},
			{
				Name:        "reserved",
				Description: "True if reserved; otherwise, false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.Reserved"),
			},
			{
				Name:        "is_xenon",
				Description: "Obsolete: Hyper-V sandbox.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.IsXenon"),
			},
			{
				Name:        "hyper_v",
				Description: "Hyper-V sandbox.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.HyperV"),
			},
			{
				Name:        "scm_site_also_stopped",
				Description: "True to stop SCM (KUDU) site when the app is stopped; otherwise, false. The default is false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.ScmSiteAlsoStopped"),
			},
			{
				Name:        "target_swap_slot",
				Description: "Specifies which deployment slot this app will swap into.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.TargetSwapSlot"),
			},
			{
				Name:        "client_affinity_enabled",
				Description: "True to enable client affinity; false to stop sending session affinity cookies, which route client requests in the same session to the same instance. Default is true.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.ClientAffinityEnabled"),
			},
			{
				Name:        "client_cert_mode",
				Description: "This composes with ClientCertEnabled setting. ClientCertEnabled: false means ClientCert is ignored. ClientCertEnabled: true and ClientCertMode: Required means ClientCert is required.ClientCertEnabled: true and ClientCertMode: Optional means ClientCert is optional or accepted. Possible values include: 'Required', 'Optional'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.ClientCertMode"),
			},
			{
				Name:        "client_cert_exclusion_paths",
				Description: "Client certificate authentication comma-separated exclusion paths.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.ClientCertExclusionPaths"),
			},
			{
				Name:        "host_names_disabled",
				Description: "True to disable the public hostnames of the app; otherwise, false. If true, the app is only accessible via API management process.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.HostNamesDisabled"),
			},
			{
				Name:        "custom_domain_verification_id",
				Description: "Unique identifier that verifies the custom domains assigned to the app. The customer will add this ID to a text record for verification.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.CustomDomainVerificationID"),
			},
			{
				Name:        "outbound_ip_addresses",
				Description: "List of IP addresses that the app uses for outbound connections (e.g. database access). Includes VIPs from tenants that site can be hosted with current settings.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.OutboundIPAddresses"),
			},
			{
				Name:        "possible_outbound_ip_addresses",
				Description: "List of IP addresses that the app uses for outbound connections (e.g. database access). Includes VIPs from all tenants except dataComponent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.PossibleOutboundIPAddresses"),
			},
			{
				Name:        "container_size",
				Description: "Size of the function container.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SiteProperties.ContainerSize"),
			},
			{
				Name:        "suspended_till",
				Description: "App suspended till in case memory-time quota is exceeded.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SiteProperties.SuspendedTill.Time"),
			},
			{
				Name:        "is_default_container",
				Description: "True if the app is a default container; otherwise, false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.IsDefaultContainer"),
			},
			{
				Name:        "default_host_name",
				Description: "Default hostname of the app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.DefaultHostName"),
			},
			{
				Name:        "https_only",
				Description: "Configures a web site to accept only https requests.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SiteProperties.HTTPSOnly"),
			},
			{
				Name:        "redundancy_mode",
				Description: "Site redundancy mode. Possible values include: 'RedundancyModeNone', 'RedundancyModeManual', 'RedundancyModeFailover', 'RedundancyModeActiveActive', 'RedundancyModeGeoRedundant'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteProperties.RedundancyMode"),
			},

			// JSON fields
			{
				Name:        "identity",
				Description: "Managed service identity.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(slotIdentity),
			},
			{
				Name:        "host_names",
				Description: "Hostnames associated with the app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.HostNames"),
			},
			{
				Name:        "enabled_host_names",
				Description: "Enabled hostnames for the app. Hostnames need to be assigned (see HostNames) AND enabled. Otherwise, the app is not served on those hostnames.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.EnabledHostNames"),
			},
			{
				Name:        "host_name_ssl_states",
				Description: "Hostname SSL states are used to manage the SSL bindings for app's hostnames.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.HostNameSslStates"),
			},
			{
				Name:        "site_config",
				Description: "Configuration of the app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.SiteConfig"),
			},
			{
				Name:        "traffic_manager_host_names",
				Description: "Azure Traffic Manager hostnames associated with the app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.TrafficManagerHostNames"),
			},
			{
				Name:        "hosting_environment_profile",
				Description: "App Service Environment to use for the app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.HostingEnvironmentProfile"),
			},
			{
				Name:        "slot_swap_status",
				Description: "Status of the last deployment slot swap operation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SiteProperties.SlotSwapStatus"),
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

type SlotInfo struct {
	SiteProperties *web.SiteProperties
	Identity       *web.ManagedServiceIdentity
	ID             *string
	Name           *string
	AppName        *string
	Kind           *string
	Location       *string
	Type           *string
	Tags           map[string]*string
}

//// LIST FUNCTION

func listAppServiceWebAppSlots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var appName, resourceGroupName string
	if h.Item != nil {
		data := h.Item.(web.Site)
		appName = *data.Name
		resourceGroupName = *data.ResourceGroup
	} else {
		return nil, nil
	}

	// Restrict the API call for other apps if the app name is specified in the query paramater
	if d.EqualsQualString("app_name") != "" && d.EqualsQualString("app_name") != appName {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_app_service_web_app_slot.listAppServiceWebAppSlots", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer
	webClient.RetryAttempts = session.RetryAttempts
	webClient.RetryDuration = session.RetryDuration

	result, err := webClient.ListSlots(ctx, resourceGroupName, appName)
	if err != nil {
		plugin.Logger(ctx).Error("azure_app_service_web_app_slot.listAppServiceWebAppSlots", "api_error", err)
		return nil, err
	}
	for _, slot := range result.Values() {
		d.StreamListItem(ctx, &SlotInfo{
			SiteProperties: slot.SiteProperties,
			Identity:       slot.Identity,
			ID:             slot.ID,
			Name:           slot.Name,
			AppName:        &appName,
			Kind:           slot.Kind,
			Location:       slot.Location,
			Type:           slot.Type,
			Tags:           slot.Tags,
		})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_app_service_web_app_slot.listAppServiceWebAppSlots", "api_pagging_error", err)
			return nil, err
		}

		for _, slot := range result.Values() {
			d.StreamListItem(ctx, slot)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAppServiceWebAppSlot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	appName := d.EqualsQuals["app_name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	slotName := d.EqualsQuals["name"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = web.AppsClient#GetSlot: Invalid input: autorest/validation: validation failed: parameter=resourceGroupName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(resourceGroup) < 1 || len(appName) < 1 || len(slotName) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_app_service_web_app_slot.getAppServiceWebAppSlot", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer
	webClient.RetryAttempts = session.RetryAttempts
	webClient.RetryDuration = session.RetryDuration

	op, err := webClient.GetSlot(ctx, resourceGroup, appName, slotName)
	if err != nil {
		plugin.Logger(ctx).Error("azure_app_service_web_app_slot.getAppServiceWebAppSlot", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return &SlotInfo{
			SiteProperties: op.SiteProperties,
			Identity:       op.Identity,
			ID:             op.ID,
			Name:           op.Name,
			AppName:        &appName,
			Kind:           op.Kind,
			Location:       op.Location,
			Type:           op.Type,
			Tags:           op.Tags,
		}, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func slotIdentity(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*SlotInfo)
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
