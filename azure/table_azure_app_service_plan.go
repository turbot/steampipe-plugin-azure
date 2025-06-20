package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureAppServicePlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_plan",
		Description: "Azure App Service Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAppServicePlan,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "serverFarms/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAppServicePlans,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "serverFarms/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the app service plan",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an app service plan uniquely",
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
				Description: "The resource type of the app service plan",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hyper_v",
				Description: "Specify whether resource is Hyper-V container app service plan",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServicePlanProperties.HyperV"),
				Default:     false,
			},
			{
				Name:        "is_spot",
				Description: "Specify whether this App Service Plan owns spot instances, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServicePlanProperties.IsSpot"),
				Default:     false,
			},
			{
				Name:        "is_xenon",
				Description: "Specify whether resource is Hyper-V container app service plan",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServicePlanProperties.IsXenon"),
				Default:     false,
			},
			{
				Name:        "maximum_elastic_worker_count",
				Description: "Maximum number of total workers allowed for this ElasticScaleEnabled App Service Plan",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServicePlanProperties.MaximumElasticWorkerCount"),
			},
			{
				Name:        "maximum_number_of_workers",
				Description: "Maximum number of instances that can be assigned to this App Service plan",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServicePlanProperties.MaximumNumberOfWorkers"),
			},
			{
				Name:        "per_site_scaling",
				Description: "Specify whether apps assigned to this App Service plan can be scaled independently",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServicePlanProperties.PerSiteScaling"),
				Default:     false,
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the App Service Plan",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServicePlanProperties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "reserved",
				Description: "Specify whether the resource is Linux app service plan, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServicePlanProperties.Reserved"),
				Default:     false,
			},
			{
				Name:        "zone_redundant",
				Description: "If true, this App Service Plan will perform availability zone balancing.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServicePlanProperties.ZoneRedundant"),
			},
			{
				Name:        "sku_capacity",
				Description: "Current number of instances assigned to the resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_family",
				Description: "Family code of the resource SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Family"),
			},
			{
				Name:        "sku_name",
				Description: "Name of the resource SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_size",
				Description: "Size specifier of the resource SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Size"),
			},
			{
				Name:        "sku_tier",
				Description: "Service tier of the resource SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "status",
				Description: "App Service plan status",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServicePlanProperties.Status").Transform(transform.ToString),
			},
			{
				Name:        "geo_region",
				Description: "Geographical location for the App Service plan",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServicePlanProperties.GeoRegion"),
			},
			{
				Name:        "elastic_scale_enabled",
				Description: "Specifies if ElasticScale is enabled for the App Service plan",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AppServicePlanProperties.ElasticScaleEnabled"),
			},
			{
				Name:        "worker_tier_name",
				Description: "Target worker tier assigned to the App Service plan",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppServicePlanProperties.WorkerTierName"),
			},
			{
				Name:        "target_worker_count",
				Description: "Scaling worker count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServicePlanProperties.TargetWorkerCount"),
			},
			{
				Name:        "target_worker_size_id",
				Description: "Scaling worker size ID",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AppServicePlanProperties.TargetWorkerSizeID"),
			},
			{
				Name:        "apps",
				Description: "Site a web app, a mobile app backend, or an API app.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServicePlanApps,
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
				Transform:   transform.FromField("AppServicePlanProperties.ResourceGroup").Transform(toLower),
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listAppServicePlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppServicePlansClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &webClient, d.Connection)

	result, err := webClient.List(ctx, types.Bool(true))
	if err != nil {
		return nil, err
	}
	for _, servicePlan := range result.Values() {
		d.StreamListItem(ctx, servicePlan)
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

		for _, servicePlan := range result.Values() {
			d.StreamListItem(ctx, servicePlan)
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

func getAppServicePlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppServicePlan")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// resourceGroupName can't be empty
	// Error: pq: rpc error: code = Unknown desc = web.AppServicePlansClient#Get: Invalid input: autorest/validation: validation failed: parameter=resourceGroupName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppServicePlansClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &webClient, d.Connection)

	op, err := webClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}
	return op, nil
}

type AppServicePlanApp struct {
	SiteProperties *web.SiteProperties
	ID             *string
	Name           *string
	Kind           *string
	Location       *string
	Type           *string
	Tags           map[string]*string
}

func getServicePlanApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	servicePlan := h.Item.(web.AppServicePlan)

	resourceGroupName := strings.Split(string(*servicePlan.ID), "/")[4]

	var apps []AppServicePlanApp

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_app_service_plan.getServicePlanApps", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	webClient := web.NewAppServicePlansClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	webClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &webClient, d.Connection)

	op, err := webClient.ListWebApps(ctx, resourceGroupName, *servicePlan.Name, "", "", "")

	if err != nil {
		plugin.Logger(ctx).Error("azure_app_service_plan.getServicePlanApps", "api_error", err)
		return nil, err
	}
	app := &AppServicePlanApp{}
	for _, data := range op.Values() {
		if data.SiteProperties != nil {
			app.SiteProperties = data.SiteProperties
		}
		if data.Name != nil {
			app.Name = data.Name
		}
		if data.ID != nil {
			app.ID = data.ID
		}
		if data.Kind != nil {
			app.Kind = data.Kind
		}
		if data.Type != nil {
			app.Type = data.Type
		}
		app.Tags = data.Tags
		apps = append(apps, *app)
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, data := range op.Values() {
			if data.SiteProperties != nil {
				app.SiteProperties = data.SiteProperties
			}
			if data.Name != nil {
				app.Name = data.Name
			}
			if data.ID != nil {
				app.ID = data.ID
			}
			if data.Kind != nil {
				app.Kind = data.Kind
			}
			if data.Type != nil {
				app.Type = data.Type
			}
			app.Tags = data.Tags
			apps = append(apps, *app)
		}
	}

	return apps, nil
}
