package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2020-02-02/insights"
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
				Description: "The friendly name that identifies the application insight",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a application insight uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "flow_type",
				Description: "Determines what kind of flow this component was created by. Possible values include: 'FlowTypeBluefield'",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.FlowType").Transform(transform.ToString),
			},
			{
				Name:        "ingestion_mode",
				Description: "Indicates the flow of the ingestion",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.IngestionMode"),
			},
			{
				Name:        "kind",
				Description: "The kind of application that this component refers to, used to customize UI. This value is a freeform string, values should typically be one of the following: web, ios, other, store, java, phone",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_network_access_for_ingestion",
				Description: "The network access type for accessing Application Insights ingestion",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.PublicNetworkAccessForIngestion"),
			},
			{
				Name:        "public_network_access_for_query",
				Description: "The network access type for accessing Application Insights query",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.PublicNetworkAccessForQuery"),
			},
			{
				Name:        "request_source",
				Description: "Describes what tool created this Application Insights component",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.RequestSource"),
			},
			{
				Name:        "retention_in_days",
				Description: "Retention period in days",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.RetentionInDays"),
			},
			{
				Name:        "type",
				Description: "The resource type of the application insight",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workspace_resource_id",
				Description: "Resource Id of the log analytics workspace which the data will be ingested to",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationInsightsComponentProperties.WorkspaceResourceID"),
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
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
