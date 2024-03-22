package azure

import (
	"context"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/maintenance/mgmt/maintenance"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureMaintenanceConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_maintenance_configuration",
		Description: "Azure Maintenance Configuration.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"resource_group", "name"}),
			Hydrate:    getMaintenanceConfiguration,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound",  "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMaintenanceConfigurations,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Fully qualified identifier of the resource.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "Name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "Gets or sets namespace of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigurationProperties.Namespace"),
			},
			{
				Name:        "visibility",
				Description: "The visibility of the configuration. The default value is 'Custom'. Possible values include: 'VisibilityCustom', 'VisibilityPublic'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigurationProperties.Visibility"),
			},
			{
				Name:        "maintenance_scope",
				Description: "The maintenanceScope of the configuration. Possible values include: 'ScopeHost', 'ScopeOSImage', 'ScopeExtension', 'ScopeInGuestPatch', 'ScopeSQLDB', 'ScopeSQLManagedInstance'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigurationProperties.MaintenanceScope"),
			},
			{
				Name:        "created_at",
				Description: "The timestamp of resource creation (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SystemData.CreatedAt").Transform(convertDateToTime),
			},
			{
				Name:        "created_by",
				Description: "The identity that created the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.CreatedBy"),
			},
			{
				Name:        "created_by_type",
				Description: "The type of identity that created the resource. Possible values include: 'CreatedByTypeUser', 'CreatedByTypeApplication', 'CreatedByTypeManagedIdentity', 'CreatedByTypeKey'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.CreatedByType"),
			},
			{
				Name:        "last_modified_at",
				Description: "The timestamp of resource last modification (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SystemData.LastModifiedAt").Transform(convertDateToTime),
			},
			{
				Name:        "last_modified_by",
				Description: "The identity that last modified the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.LastModifiedBy"),
			},
			{
				Name:        "last_modified_by_type",
				Description: "The type of identity that last modified the resource. Possible values include: 'CreatedByTypeUser', 'CreatedByTypeApplication', 'CreatedByTypeManagedIdentity', 'CreatedByTypeKey'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemData.LastModifiedByType"),
			},
			{
				Name:        "extension_properties",
				Description: "Gets or sets extensionProperties of the maintenanceConfiguration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConfigurationProperties.ExtensionProperties"),
			},
			{
				Name:        "window",
				Description: "Definition of a MaintenanceWindow.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConfigurationProperties.Window"),
			},
			{
				Name:        "system_data",
				Description: "Azure Resource Manager metadata containing createdBy and modifiedBy information.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SystemData").Transform(extractConfigurationSystemData),
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

//// LIST FUNCTION

func listMaintenanceConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_maintenance_configuration.listMaintenanceConfigurations", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := maintenance.NewConfigurationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// The API doesn't support pagination
	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_maintenance_configuration.listMaintenanceConfigurations", "api_error", err)
		return nil, err
	}
	for _, res := range *result.Value {
		d.StreamListItem(ctx, res)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMaintenanceConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var resourceGroup, name string
	if h.Item != nil {
		name = *h.Item.(maintenance.Configuration).Name
	} else {
		name = d.EqualsQualString("name")
		resourceGroup = d.EqualsQualString("resource_group")
	}

	// check if name is empty
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_maintenance_configuration.getMaintenanceConfiguration", "session_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID

	client := maintenance.NewConfigurationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_maintenance_configuration.getMaintenanceConfiguration", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func extractConfigurationSystemData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	conf := d.HydrateItem.(maintenance.Configuration)
	if conf.SystemData != nil {
		return structToMap(reflect.ValueOf(*conf.SystemData)), nil
	}

	return nil, nil
}
