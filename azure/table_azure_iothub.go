package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureIotHub(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_iothub",
		Description: "Azure Iot Hub",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getIotHub,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listIotHubs,
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
				Description: "The iot hub state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.State"),
			},
			{
				Name:        "provisioning_state",
				Description: "The iot hub provisioning state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "comments",
				Description: "Iot hub comments.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Comments"),
			},
			{
				Name:        "enable_file_upload_notifications",
				Description: "Indicates if file upload notifications are enabled for the iot hub.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.EnableFileUploadNotifications"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "features",
				Description: "The capabilities and features enabled for the iot hub. Possible values include: 'None', 'DeviceManagement'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Features"),
			},
			{
				Name:        "host_name",
				Description: "The name of the host.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HostName"),
			},
			{
				Name:        "min_tls_version",
				Description: "Specifies the minimum TLS version to support for this iot hub.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinTLSVersion"),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether requests from public network are allowed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "sku_capacity",
				Description: "Iot hub SKU capacity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_name",
				Description: "Iot hub SKU name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "Iot hub SKU tier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "authorization_policies",
				Description: "The shared access policies you can use to secure a connection to the iot hub.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AuthorizationPolicies"),
			},
			{
				Name:        "cloud_to_device",
				Description: "CloudToDevice properties of the iot hub.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CloudToDevice"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the iot hub.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listIotHubDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "event_hub_endpoints",
				Description: "The event hub-compatible endpoint properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.EventHubEndpoints"),
			},
			{
				Name:        "ip_filter_rules",
				Description: "The IP filter rules of the iot hub.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.IPFilterRules"),
			},
			{
				Name:        "locations",
				Description: "Primary and secondary location for iot hub.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Locations"),
			},
			{
				Name:        "messaging_endpoints",
				Description: "The messaging endpoint properties for the file upload notification queue.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.MessagingEndpoints"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "Private endpoint connections created on this iot hub.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PrivateEndpointConnections"),
			},
			{
				Name:        "routing",
				Description: "Routing properties of the iot hub.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Routing"),
			},
			{
				Name:        "storage_endpoints",
				Description: "The list of azure storage endpoints where you can upload files.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.StorageEndpoints"),
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

func listIotHubs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	iotHubClient := devices.NewIotHubResourceClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	iotHubClient.Authorizer = session.Authorizer
	result, err := iotHubClient.ListBySubscription(ctx)
	if err != nil {
		return nil, err
	}
	for _, iotHubDescription := range result.Values() {
		d.StreamListItem(ctx, iotHubDescription)
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
		for _, iotHubDescription := range result.Values() {
			d.StreamListItem(ctx, iotHubDescription)
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

func getIotHub(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIotHub")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	iotHubClient := devices.NewIotHubResourceClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	iotHubClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := iotHubClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listIotHubDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listIotHubDiagnosticSettings")
	id := *h.Item.(devices.IotHubDescription).ID

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
