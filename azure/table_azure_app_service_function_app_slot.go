package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

//// TABLE DEFINITION

// The slots of the apps azure_app_service_function_app lists. The list and get hydrates are the
// web app slot table's: a slot is fetched the same way whatever kind of app it belongs to.
func tableAzureAppServiceFunctionAppSlot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_function_app_slot",
		Description: "Azure App Service Function App Slot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "app_name", "resource_group"}),
			Hydrate:    getAppServiceWebAppSlot,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "sites/slots/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAppServiceFunctionApps,
			Hydrate:       listAppServiceWebAppSlots,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "sites/slots/read",
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "app_name",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getConfigurationSlot,
				Tags: map[string]string{
					"service": "Microsoft.Web",
					"action":  "sites/slots/config/read",
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Resource Name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
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
			{
				Name:        "site_config_resource",
				Description: "Configuration of an app, such as platform version and bitness, default documents, virtual applications, Always On, etc.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConfigurationSlot,
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
