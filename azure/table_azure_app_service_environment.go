package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAppServiceEnvironment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_environment",
		Description: "Azure App Service Environment",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAppServiceEnvironment,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAppServiceEnvironments,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the app service environment.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an app service environment uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "kind",
				Description: "Contains the kind of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the app service environment. Possible values include: 'Succeeded', 'Failed', 'Canceled', 'InProgress', 'Deleting'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "The resource type of the app service environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allowed_multi_sizes",
				Description: "List of comma separated strings describing which VM sizes are allowed for front-ends.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.AllowedMultiSizes"),
			},
			{
				Name:        "allowed_worker_sizes",
				Description: "List of comma separated strings describing which VM sizes are allowed for workers.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.AllowedWorkerSizes"),
			},
			{
				Name:        "api_management_account_id",
				Description: "API management account associated with the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.APIManagementAccountID"),
			},
			{
				Name:        "database_edition",
				Description: "Edition of the metadata database for the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.DatabaseEdition"),
			},
			{
				Name:        "database_service_objective",
				Description: "Service objective of the metadata database for the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.DatabaseServiceObjective"),
			},
			{
				Name:        "default_front_end_scale_factor",
				Description: "Default Scale Factor for FrontEnds.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.DefaultFrontEndScaleFactor"),
			},
			{
				Name:        "dns_suffix",
				Description: "DNS suffix of the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.DNSSuffix"),
			},
			{
				Name:        "dynamic_cache_enabled",
				Description: "Indicates whether the dynamic cache is enabled or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.DynamicCacheEnabled"),
				Default:     false,
			},
			{
				Name:        "environment_status",
				Description: "Detailed message about with results of the last check of the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.EnvironmentStatus"),
			},
			{
				Name:        "front_end_scale_factor",
				Description: "Scale factor for front-ends.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.FrontEndScaleFactor"),
			},
			{
				Name:        "has_linux_workers",
				Description: "Indicates whether an ASE has linux workers or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.HasLinuxWorkers"),
				Default:     false,
			},
			{
				Name:        "internal_load_balancing_mode",
				Description: "Specifies which endpoints to serve internally in the virtual network for the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.InternalLoadBalancingMode").Transform(transform.ToString),
			},
			{
				Name:        "ip_ssl_address_count",
				Description: "Number of IP SSL addresses reserved for the app service environment.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.IpsslAddressCount"),
			},
			{
				Name:        "is_healthy_environment",
				Description: "Indicates whether the app service environment is healthy.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.EnvironmentIsHealthy"),
				Default:     false,
			},
			{
				Name:        "last_action",
				Description: "Last deployment action on the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.LastAction"),
			},
			{
				Name:        "last_action_result",
				Description: "Result of the last deployment action on the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.LastActionResult"),
			},
			{
				Name:        "maximum_number_of_machines",
				Description: "Maximum number of VMs in the app service environment.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.MaximumNumberOfMachines"),
			},
			{
				Name:        "multi_role_count",
				Description: "Number of front-end instances.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.MultiRoleCount"),
			},
			{
				Name:        "multi_size",
				Description: "Front-end VM size.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.MultiSize"),
			},
			{
				Name:        "ssl_cert_key_vault_id",
				Description: "Key vault ID for ILB app service environment default SSL certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.SslCertKeyVaultID"),
			},
			{
				Name:        "ssl_cert_key_vault_secret_name",
				Description: "Key vault secret name for ILB app service environment default SSL certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.SslCertKeyVaultSecretName"),
			},
			{
				Name:        "status",
				Description: "Current status of the app service environment. Possible values include: 'Preparing', 'Ready', 'Scaling', 'Deleting'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.Status").Transform(transform.ToString),
			},
			{
				Name:        "suspended",
				Description: "Indicates whether the app service environment is suspended or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.Suspended"),
				Default:     false,
			},
			{
				Name:        "upgrade_domains",
				Description: "Number of upgrade domains of the app service environment.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.UpgradeDomains"),
			},
			{
				Name:        "vnet_name",
				Description: "Name of the virtual network for the app service environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.VnetName"),
			},
			{
				Name:        "vnet_resource_group_name",
				Description: "Name of the resource group where the virtual network is created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.VnetResourceGroupName"),
			},
			{
				Name:        "vnet_subnet_name",
				Description: "Name of the subnet of the virtual network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.VnetSubnetName"),
			},
			{
				Name:        "cluster_settings",
				Description: "Custom settings for changing the behavior of the app service environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppServiceEnvironment.ClusterSettings"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the app service environment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAppServiceEnvironmentDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "environment_capacities",
				Description: "Current total, used, and available worker capacities.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppServiceEnvironment.EnvironmentCapacities"),
			},
			{
				Name:        "network_access_control_list",
				Description: " Access control list for controlling traffic to the app service environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppServiceEnvironment.NetworkAccessControlList"),
			},
			{
				Name:        "user_whitelisted_ip_ranges",
				Description: "User added ip ranges to whitelist on ASE db.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppServiceEnvironment.UserWhitelistedIPRanges"),
			},
			{
				Name:        "vip_mappings",
				Description: "Description of IP SSL mapping for the app service environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppServiceEnvironment.VipMappings"),
			},
			{
				Name:        "virtual_network",
				Description: "Description of the virtual network.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractAppServiceEnvironmentVirtualNetwork),
			},
			{
				Name:        "worker_pools",
				Description: "Description of worker pools with worker size IDs, VM sizes, and number of workers in each pool.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppServiceEnvironment.WorkerPools"),
			},

			// Standard columns
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
				Transform:   transform.FromField("AppServiceEnvironment.ResourceGroup").Transform(toLower),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.SubscriptionID"),
			},
		},
	}
}

//// LIST FUNCTION

func listAppServiceEnvironments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppServiceEnvironmentsClient(subscriptionID)
	webClient.Authorizer = session.Authorizer

	result, err := webClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, environment := range result.Values() {
		d.StreamListItem(ctx, environment)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, environment := range result.Values() {
			d.StreamListItem(ctx, environment)
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAppServiceEnvironment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceEnvironment")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// resourceGroupName can't be empty
	// Error: pq: rpc error: code = Unknown desc = web.AppServiceEnvironmentsClient#Get: Invalid input: autorest/validation: validation failed: parameter=resourceGroupName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppServiceEnvironmentsClient(subscriptionID)
	webClient.Authorizer = session.Authorizer

	op, err := webClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func listAppServiceEnvironmentDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAppServiceEnvironmentDiagnosticSettings")
	id := *h.Item.(web.AppServiceEnvironmentResource).ID

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

	// If we return the API response directly, the output does not provide
	// all the contents of DiagnosticSettings
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

//// TRANSFORM FUNCTION

// If we return the API response directly, the output does not provide
// all the properties of VirtualNetwork
func extractAppServiceEnvironmentVirtualNetwork(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	env := d.HydrateItem.(web.AppServiceEnvironmentResource)
	properties := make(map[string]interface{})

	if env.AppServiceEnvironment != nil && env.AppServiceEnvironment.VirtualNetwork != nil {
		if env.AppServiceEnvironment.VirtualNetwork.ID != nil {
			properties["id"] = env.AppServiceEnvironment.VirtualNetwork.ID
		}
		if env.AppServiceEnvironment.VirtualNetwork.Name != nil {
			properties["name"] = env.AppServiceEnvironment.VirtualNetwork.Name
		}
		if env.AppServiceEnvironment.VirtualNetwork.Type != nil {
			properties["type"] = env.AppServiceEnvironment.VirtualNetwork.Type
		}
		if env.AppServiceEnvironment.VirtualNetwork.ID != nil {
			properties["subnet"] = env.AppServiceEnvironment.VirtualNetwork.Subnet
		}
	}

	return properties, nil
}
