package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/appplatform/mgmt/2020-07-01/appplatform"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureSpringCloudApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_spring_cloud_app",
		Description: "Azure Spring Cloud App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group", "service_name"}),
			Hydrate:    getSpringCloudApp,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listSpringCloudServicesBySubscription,
			Hydrate:       listSpringCloudApps,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "service_name",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_group",
					Require: plugin.Optional,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource Id for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The name of the service where the app is belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").TransformP(extractPropertyByIndexSplitByBackSlash, 8),
			},
			{
				Name:        "location",
				Description: "The GEO location of the application, always the same with its parent resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public",
				Description: "Indicates whether the App exposes public endpoint.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.Public"),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the App. Possible values include: 'Succeeded', 'Failed', 'Creating', 'Updating'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "created_time",
				Description: "Date time when the resource is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.CreatedTime.Time"),
			},
			{
				Name:        "url",
				Description: "URL of the App.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.URL"),
			},
			{
				Name:        "active_deployment_name",
				Description: "Name of the active deployment of the App.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ActiveDeploymentName"),
			},
			{
				Name:        "fqdn",
				Description: "Fully qualified dns Name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Fqdn"),
			},
			{
				Name:        "https_only",
				Description: "Indicate if only https is allowed.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.HTTPSOnly"),
			},
			{
				Name:        "persistent_disk",
				Description: "Persistent disk settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PersistentDisk"),
			},
			{
				Name:        "temporary_disk",
				Description: "Temporary disk settings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.TemporaryDisk"),
			},
			{
				Name:        "identity",
				Description: "The Managed Identity type of the app resource.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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

// As of 15th January 2024, Steampipe does not support ParentHydrate for more than two levels.
// In this hierarchy, an App is a sub-element of Service, and Service is a sub-element of Resource Group (Spring Cloud App > Spring Cloud Service > Resource Group).
// Since azure_resource_group is the parent of azure_spring_cloud_service,
// azure_spring_cloud_service cannot be the parent table for azure_spring_cloud_app.

// Implementation Steps:
/*
1. Create a function to list all Azure Spring Cloud Services by subscription and streamlist the results in this table.
2. Utilize this function as the parent hydrate for the table azure_spring_cloud_app.
3. Designate 'service_name' and 'resource_group_name' as optional qualifiers for this table.
4. Implement checks to minimize API calls when specific service_name or resource_group are provided.
*/


//// LIST FUNCTION

func listSpringCloudApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_spring_cloud_app.listSpringCloudApps", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	// Get the details of the resource group
	service := h.Item.(appplatform.ServiceResource)
	if service.ID == nil {
		return nil, nil
	}

	resourGroup := strings.Split(*service.ID, "/")[4]

	// Limit the API calls for a given value of service name or resource group name.
	if d.EqualsQualString("resource_group") != "" {
		if resourGroup != d.EqualsQualString("resource_group") {
			return nil, nil
		}
	}

	if d.EqualsQualString("service_name") != "" {
		if *service.Name != d.EqualsQualString("service_name") {
			return nil, nil
		}
	}

	client := appplatform.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx, resourGroup, *service.Name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_spring_cloud_app.listSpringCloudApps", "api_error", err)
		return nil, err
	}

	for _, app := range result.Values() {
		d.StreamListItem(ctx, app)

		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("lazure_spring_cloud_app.listSpringCloudServices", "list_paging", err)
			return nil, err
		}
		for _, app := range result.Values() {
			d.StreamListItem(ctx, app)

			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// PARENT HYDRATE FUNCTIONS

// As of 15th January 2024, Steampipe does not support three-level deep ParentHydrate.
// In the hierarchy, the App is a child of Service, which in turn is a child of Resource Group (Structured as: Spring Cloud App > Spring Cloud Service > Resource Group).
// Therefore, azure_spring_cloud_service cannot serve as the parent table for azure_spring_cloud_app.
// Additionally, the azure_spring_cloud_service table lists services by resource groups.
func listSpringCloudServicesBySubscription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("lazure_spring_cloud_app.listSpringCloudServicesBySubscription", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := appplatform.NewServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("lazure_spring_cloud_app.listSpringCloudServicesBySubscription", "api_error", err)
		return nil, err
	}
	for _, service := range result.Values() {
		d.StreamListItem(ctx, service)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("lazure_spring_cloud_app.listSpringCloudServicesBySubscription", "paging_error", err)
			return nil, err
		}
		for _, service := range result.Values() {
			d.StreamListItem(ctx, service)
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSpringCloudApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	resourceGroup := d.EqualsQualString("resource_group")
	serviceName := d.EqualsQualString("service_name")
	name := d.EqualsQualString("name")

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" || serviceName == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("lazure_spring_cloud_app.getSpringCloudApp", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := appplatform.NewAppsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	app, err := client.Get(ctx, resourceGroup, serviceName, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("lazure_spring_cloud_app.getSpringCloudApp", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if app.ID != nil {
		return app, nil
	}

	return nil, nil
}