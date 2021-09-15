package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/machinelearningservices/mgmt/2020-02-18-preview/machinelearningservices"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureMachineLearningWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_machine_learning_workspace",
		Description: "Azure Machine Learning Workspace",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getMachineLearningWorkspace,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listMachineLearningWorkspaces,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "friendly_name",
				Description: "The friendly name for this workspace. This name in mutable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.FriendlyName"),
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The current deployment state of workspace resource, The provisioningState is to indicate states for resource provisioning. Possible values include: 'Unknown', 'Updating', 'Creating', 'Deleting', 'Succeeded', 'Failed', 'Canceled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ProvisioningState"),
			},
			{
				Name:        "creation_time",
				Description: "The creation time for this workspace resource.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("WorkspaceProperties.CreationTime").Transform(convertDateToTime),
			},
			{
				Name:        "workspace_id",
				Description: "The immutable id associated with this workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.WorkspaceID"),
			},
			{
				Name:        "application_insights",
				Description: "ARM id of the application insights associated with this workspace. This cannot be changed once the workspace has been created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ApplicationInsights"),
			},
			{
				Name:        "container_registry",
				Description: "ARM id of the container registry associated with this workspace. This cannot be changed once the workspace has been created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ContainerRegistry"),
			},
			{
				Name:        "description",
				Description: "The description of this workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.Description"),
			},
			{
				Name:        "discovery_url",
				Description: "ARM id of the container registry associated with this workspace. This cannot be changed once the workspace has been created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.DiscoveryURL"),
			},
			{
				Name:        "hbi_workspace",
				Description: "The flag to signal HBI data in the workspace and reduce diagnostic data collected by the service.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("WorkspaceProperties.HbiWorkspace"),
			},
			{
				Name:        "key_vault",
				Description: "ARM id of the key vault associated with this workspace, This cannot be changed once the workspace has been created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.KeyVault"),
			},
			{
				Name:        "location",
				Description: "The location of the resource. This cannot be changed after the resource is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_provisioned_resource_group",
				Description: "The name of the managed resource group created by workspace RP in customer subscription if the workspace is CMK workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ServiceProvisionedResourceGroup"),
			},
			{
				Name:        "sku_name",
				Description: "Name of the sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "Tier of the sku like Basic or Enterprise.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "storage_account",
				Description: "ARM id of the storage account associated with this workspace. This cannot be changed once the workspace has been created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.StorageAccount"),
			},
			{
				Name:        "studio_endpoint",
				Description: "The regional endpoint for the machine learning studio service which hosts this workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the n workspace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMachineLearningDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "The encryption settings of Azure ML workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.Encryption"),
			},
			{
				Name:        "identity",
				Description: "The identity of the resource.",
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

func listMachineLearningWorkspaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	worspaceClient := machinelearningservices.NewWorkspacesClient(subscriptionID)
	worspaceClient.Authorizer = session.Authorizer

	result, err := worspaceClient.ListBySubscription(ctx, "")
	if err != nil {
		logger.Error("listMachineLearningWorkspaces", "list", err)
		return nil, err
	}
	for _, workspace := range result.Values() {
		d.StreamListItem(ctx, workspace)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, workspace := range result.Values() {
			d.StreamListItem(ctx, workspace)
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getMachineLearningWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMachineLearningWorkspace")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	workspaceClient := machinelearningservices.NewWorkspacesClient(subscriptionID)
	workspaceClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	workspace, err := workspaceClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getMachineLearningWorkspace", "get", err)
		return nil, err
	}

	return workspace, nil
}

func listMachineLearningDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMachineLearningDiagnosticSettings")
	id := *h.Item.(machinelearningservices.Workspace).ID

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
		plugin.Logger(ctx).Error("listMachineLearningDiagnosticSettings", "Error", err)
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
