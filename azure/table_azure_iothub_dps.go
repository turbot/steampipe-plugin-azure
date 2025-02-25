package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/provisioningservices/mgmt/iothub"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureIotHubDps(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_iothub_dps",
		Description: "Azure Iot Hub Dps",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getIotHubDps,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIotHubDpses,
		},
		Columns: azureColumns([]*plugin.Column{
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
				Description: "Current state of the provisioning service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.State"),
			},
			{
				Name:        "provisioning_state",
				Description: "The ARM provisioning state of the provisioning service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allocation_policy",
				Description: "Allocation policy to be used by this provisioning service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AllocationPolicy"),
			},
			{
				Name:        "device_provisioning_host_name",
				Description: "Device endpoint for this provisioning service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DeviceProvisioningHostName"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id_scope",
				Description: "Unique identifier of this provisioning service..",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.IDScope"),
			},
			{
				Name:        "service_operations_host_name",
				Description: "Service endpoint for provisioning service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ServiceOperationsHostName"),
			},
			{
				Name:        "sku_capacity",
				Description: "Iot dps SKU capacity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_name",
				Description: "Iot dps SKU name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Iot dps SKU tier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "authorization_policies",
				Description: "List of authorization keys for a provisioning service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AuthorizationPolicies"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the iot dps.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listIotDpsDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "iot_hubs",
				Description: "List of IoT hubs associated with this provisioning service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.IotHubs"),
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

func listIotHubDpses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	iotDpsClient := iothub.NewIotDpsResourceClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	iotDpsClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &iotDpsClient, d.Connection)

	result, err := iotDpsClient.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listIotHubDpses", "ListBySubscription", err)
		return nil, err
	}
	for _, provisioningServiceDescription := range result.Values() {
		d.StreamListItem(ctx, provisioningServiceDescription)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listIotHubDpses", "ListBySubscription_pagination", err)
			return nil, err
		}
		for _, provisioningServiceDescription := range result.Values() {
			d.StreamListItem(ctx, provisioningServiceDescription)
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getIotHubDps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIotHubDps")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	iotDpsClient := iothub.NewIotDpsResourceClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	iotDpsClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &iotDpsClient, d.Connection)

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := iotDpsClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getIotHubDps", "Get", err)
		return nil, err
	}

	return op, nil
}

func listIotDpsDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listIotDpsDiagnosticSettings")
	id := *h.Item.(iothub.ProvisioningServiceDescription).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.List(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("listIotDpsDiagnosticSettings", "List", err)
		return nil, err
	}

	// If we return the API response directly, the output does not provide all
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
