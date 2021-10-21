package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureLogicAppWorkflow(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_logic_app_workflow",
		Description: "Azure Logic App Workflow",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getLogicAppWorkflow,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listLogicAppWorkflows,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "state",
				Description: "The state of the workflow.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkflowProperties.State"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the workflow.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkflowProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "access_endpoint",
				Description: "The access endpoint of the workflow.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkflowProperties.AccessEndpoint"),
			},
			{
				Name:        "created_time",
				Description: "The time when workflow was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("WorkflowProperties.CreatedTime").Transform(convertDateToTime),
			},
			{
				Name:        "changed_time",
				Description: "Specifies the time, the workflow was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("WorkflowProperties.ChangedTime").Transform(convertDateToTime),
			},
			{
				Name:        "sku_name",
				Description: "The sku name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkflowProperties.Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "version",
				Description: "Version of the workflow.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkflowProperties.Version"),
			},
			{
				Name:        "access_control",
				Description: "The access control configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkflowProperties.AccessControl"),
			},
			{
				Name:        "definition",
				Description: "The workflow defination.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkflowProperties.Definition"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the workflow.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listLogicAppWorkflowDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "endpoints_configuration",
				Description: "The endpoints configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkflowProperties.EndpointsConfiguration"),
			},
			{
				Name:        "integration_account",
				Description: "The integration account of the workflow.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkflowProperties.IntegrationAccount"),
			},
			{
				Name:        "integration_service_environment",
				Description: "The integration service environment of the workflow.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkflowProperties.IntegrationServiceEnvironment"),
			},
			{
				Name:        "parameters",
				Description: "The workflow parameters.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkflowProperties.Parameters"),
			},
			{
				Name:        "sku_plan",
				Description: "The sku Plan.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkflowProperties.Sku.Plan"),
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

			// Azure standard column
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listLogicAppWorkflows(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	workflowClient := logic.NewWorkflowsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	workflowClient.Authorizer = session.Authorizer
	result, err := workflowClient.ListBySubscription(ctx, nil, "")
	if err != nil {
		return nil, err
	}
	for _, workflow := range result.Values() {
		d.StreamListItem(ctx, workflow)
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
		for _, workflow := range result.Values() {
			d.StreamListItem(ctx, workflow)
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

func getLogicAppWorkflow(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLogicAppWorkflow")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	workflowClient := logic.NewWorkflowsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	workflowClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := workflowClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listLogicAppWorkflowDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listLogicAppWorkflowDiagnosticSettings")
	id := *h.Item.(logic.Workflow).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of DiagnosticSettings
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
