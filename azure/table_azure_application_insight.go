package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/appinsights/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureApplicationInsight(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_application_insight",
		Description: "Azure Application Insight",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getApplicationInsight,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listApplicationInsights,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the application insight.",
			},
			{
				Name:        "id",
				Description: "Contains id to identify the application insight uniquely.",
				Transform:   transform.FromGo(),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "app_id",
				Description: "Application insights unique id for your Application.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.AppID"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_string",
				Description: "Application Insights component connection string.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.ConnectionString"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "Creation date for the Application Insights component.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.CreationDate").Transform(convertDateToTime),
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "disable_ip_masking",
				Description: "Disable IP masking.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.DisableIPMasking"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "disable_local_auth",
				Description: "Disable Non-AAD based Auth.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.DisableLocalAuth"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "force_customer_storage_for_profiler",
				Description: "Force users to create their own storage account for profiler and debugger.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.ForceCustomerStorageForProfiler"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "immediate_purge_data_on_30_days",
				Description: "Purge data immediately after 30 days.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.ImmediatePurgeDataOn30Days"),
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "instrumentation_key",
				Description: "Application Insights Instrumentation key.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.InstrumentationKey"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The kind of application that this component refers to, used to customize UI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Current state of this component.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.ProvisioningState"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "request_source",
				Description: "Describes what tool created this Application Insights component.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.RequestSource"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retention_in_days",
				Description: "Retention period in days.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.RetentionInDays"),
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "sampling_percentage",
				Description: "Percentage of the data produced by the application being monitored that is being sampled for Application Insights telemetry.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.SamplingPercentage"),
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "tenant_id",
				Description: "Azure Tenant ID.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.TenantID"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the application insight.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workspace_resource_id",
				Description: "Resource Id of the log analytics workspace to which the data will be ingested.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.WorkspaceResourceID"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_network_access_for_ingestion",
				Description: "The network access type for accessing Application Insights ingestion.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.PublicNetworkAccessForIngestion"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_network_access_for_query",
				Description: "The network access type for accessing Application Insights query.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.PublicNetworkAccessForQuery"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_type",
				Description: "Type of application being monitored.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.ApplicationType"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "flow_type",
				Description: "Determines what kind of flow this component was created by.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.FlowType"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ingestion_mode",
				Description: "Indicates the flow of the ingestion.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.IngestionMode"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_link_scoped_resources",
				Description: "List of linked private link scope resources.",
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.PrivateLinkScopedResources"),
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Name"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Transform:   transform.FromField("ID").Transform(idToAkas),
				Type:        proto.ColumnType_JSON,
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Transform:   transform.FromField("Location").Transform(toLower),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listApplicationInsights(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_application_insight.listApplicationInsights", "connection_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	applicationInsightClient := insights.NewComponentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	applicationInsightClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &applicationInsightClient, d.Connection)

	result, err := applicationInsightClient.List(ctx)
	if err != nil {
		logger.Error("azure_application_insight.listApplicationInsights", "api_error", err)
		return nil, err
	}

	for _, applicationInsight := range result.Values() {
		d.StreamListItem(ctx, applicationInsight)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			logger.Error("azure_application_insight.listApplicationInsights", "paging_error", err)
			return nil, err
		}

		for _, applicationInsight := range result.Values() {
			d.StreamListItem(ctx, applicationInsight)
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

func getApplicationInsight(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_application_insight.getApplicationInsight", "connection_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	applicationInsightClient := insights.NewComponentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	applicationInsightClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &applicationInsightClient, d.Connection)

	op, err := applicationInsightClient.Get(ctx, resourceGroup, name)
	if err != nil {
		logger.Error("azure_application_insight.getApplicationInsight", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
