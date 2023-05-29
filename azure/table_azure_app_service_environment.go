package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAppServiceEnvironment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_environment",
		Description: "Azure App Service Environment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAppServiceEnvironment,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAppServiceEnvironments,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the app service environment",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an app service environment uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "kind",
				Description: "Contains the kind of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the app service environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Current status of the App Service Environment",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.Status").Transform(transform.ToString),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the App Service Environment",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "default_front_end_scale_factor",
				Description: "Default Scale Factor for FrontEnds",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.DefaultFrontEndScaleFactor"),
			},
			{
				Name:        "dynamic_cache_enabled",
				Description: "Indicates whether the dynamic cache is enabled or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.EnableAcceleratedNetworking"),
				Default:     false,
			},
			{
				Name:        "front_end_scale_factor",
				Description: "Scale factor for front-ends",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServiceEnvironment.FrontEndScaleFactor"),
			},
			{
				Name:        "has_linux_workers",
				Description: "Indicates whether an ASE has linux workers or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.HasLinuxWorkers"),
				Default:     false,
			},
			{
				Name:        "internal_load_balancing_mode",
				Description: "Specifies which endpoints to serve internally in the Virtual Network for the App Service Environment",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.InternalLoadBalancingMode").Transform(transform.ToString),
			},
			{
				Name:        "is_healthy_environment",
				Description: "Indicates whether the App Service Environment is healthy",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.EnvironmentIsHealthy"),
				Default:     false,
			},
			{
				Name:        "suspended",
				Description: "Indicates whether the App Service Environment is suspended or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServiceEnvironment.Suspended"),
				Default:     false,
			},
			{
				Name:        "vnet_name",
				Description: "Name of the Virtual Network for the App Service Environment",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.VnetName"),
			},
			{
				Name:        "vnet_resource_group_name",
				Description: "Name of the resource group where the virtual network is created",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.VnetResourceGroupName"),
			},
			{
				Name:        "vnet_subnet_name",
				Description: "Name of the subnet of the virtual network",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServiceEnvironment.VnetSubnetName"),
			},
			{
				Name:        "cluster_settings",
				Description: "Custom settings for changing the behavior of the App Service Environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppServiceEnvironment.ClusterSettings"),
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
				Transform:   transform.FromField("AppServiceEnvironment.ResourceGroup").Transform(toLower),
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listAppServiceEnvironments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppServiceEnvironmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	result, err := webClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, environment := range result.Values() {
		d.StreamListItem(ctx, environment)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, environment := range result.Values() {
			d.StreamListItem(ctx, environment)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getAppServiceEnvironment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServiceEnvironment")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

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

	webClient := web.NewAppServiceEnvironmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	op, err := webClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}
	return op, nil
}
