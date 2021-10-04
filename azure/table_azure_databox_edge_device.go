package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2019-07-01/databoxedge"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureDataBoxEdgeDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_databox_edge_device",
		Description: "Azure Databox Edge Device",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getDataLboxEdgeDevice,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listDataboxEdgeDevices,
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
				Name:        "friendly_name",
				Description: "The Data Box Edge/Gateway device name.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.FriendlyName"),
			},
			{
				Name:        "data_box_edge_device_status",
				Description: "The status of the Data Box Edge/Gateway device. Possible values include: 'ReadyToSetup', 'Online', 'Offline', 'NeedsAttention', 'Disconnected', 'PartiallyDisconnected', 'Maintenance'.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.DataBoxEdgeDeviceStatus"),
			},
			{
				Name:        "Description",
				Description: "he Description of the Data Box Edge/Gateway device.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.Description"),
			},
			{
				Name:        "model_description",
				Description: "The description of the Data Box Edge/Gateway device model.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.ModelDescription"),
			},
			{
				Name:        "device_type",
				Description: "The type of the Data Box Edge/Gateway device. Possible values include: 'DataBoxEdgeDevice'.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.DeviceType"),
			},
			{
				Name:        "culture",
				Description: "The Data Box Edge/Gateway device culture.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.Culture"),
			},
			{
				Name:        "device_model",
				Description: "The Data Box Edge/Gateway device model.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.DeviceModel"),
			},
			{
				Name:        "device_software_version",
				Description: "The Data Box Edge/Gateway device software version.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.DeviceSoftwareVersion"),
			},
			{
				Name:        "device_local_capacity",
				Description: "The Data Box Edge/Gateway device local capacity in MB.",
				Type:        proto.ColumnType_INT,
				Transform: transform.FromField("DeviceProperties.DeviceLocalCapacity"),
			},
			{
				Name:        "time_zone",
				Description: "The Data Box Edge/Gateway device timezone.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.TimeZone"),
			},
			{
				Name:        "device_hcs_version",
				Description: "The device software version number of the device (eg: 1.2.18105.6).",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.DeviceHcsVersion"),
			},
			{
				Name:        "serial_number",
				Description: "The Serial Number of Data Box Edge/Gateway device.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DeviceProperties.SerialNumber"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The etag for the devices.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the device. This is a supported and registered Azure geographical region (for example, West US, East US, or Southeast Asia). The geographical region of a device cannot be changed once it is created, but if an identical geographical region is specified on update, the request will succeed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku_name",
				Description: "SKU name of the resource. Possible values include: 'Gateway', 'Edge'.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "The SKU tier. This is based on the SKU name. Possible values include: 'Standard'.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("Sku.Tier"),
			},
			{
				Name:        "node_count",
				Description: "The number of nodes in the cluster.",
				Type:        proto.ColumnType_INT,
				Transform: transform.FromField("DeviceProperties.NodeCount"),
			},
			{
				Name:        "configured_role_types",
				Description: "Type of compute roles configured.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField("DeviceProperties.ConfiguredRoleTypes"),
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

func listDataboxEdgeDevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDataboxEdgeDevices")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	deviceClient := databoxedge.NewDevicesClient(subscriptionID)
	deviceClient.Authorizer = session.Authorizer

	result, err := deviceClient.ListBySubscription(ctx, "")
	if err != nil {
		plugin.Logger(ctx).Error("listDataboxEdgeDevices", "ListBySubscription", err)
		return nil, err
	}
	for _, account := range result.Values() {
		d.StreamListItem(ctx, account)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listDataboxEdgeDevices", "ListBySubscription_pagination", err)
			return nil, err
		}
		for _, account := range result.Values() {
			d.StreamListItem(ctx, account)
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDataLboxEdgeDevice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataLboxEdgeDevice")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	deviceClient := databoxedge.NewDevicesClient(subscriptionID)
	deviceClient.Authorizer = session.Authorizer

	var name, resourceGroup string
	if h.Item != nil {
		data := h.Item.(account.DataLakeStoreAccountBasic)
		splitID := strings.Split(*data.ID, "/")
		name = *data.Name
		resourceGroup = splitID[4]
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroup = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := deviceClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
