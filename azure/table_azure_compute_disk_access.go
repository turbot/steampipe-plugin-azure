package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureComputeDiskAccess(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_disk_access",
		Description: "Azure Compute Disk Access",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAzureComputeDiskAccess,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureComputeDiskAccesses,
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
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The disk access resource provisioning state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time when the disk access was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated").Transform(convertDateToTime),
			},
			{
				Name:        "private_endpoints_id",
				Description: "The private endpoints ids.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointsID"),
			},
			{
				Name:        "private_endpoint_connections_id",
				Description: "The private endpoint connections ids.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointConnectionsID"),
			},
			{
				Name:        "private_endpoint_connections_name",
				Description: "The private endpoint connections names.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_endpoint_connections_type",
				Description: "The private endpoint connections types.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The private endpoint connections details.",
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

type diskAccesssInfo struct {
	// ID - READ-ONLY; Resource Id
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type
	Type *string `json:"type,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags
	Tags map[string]*string `json:"tags"`
	// ProvisioningState - READ-ONLY; The disk access resource provisioning state.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// TimeCreated - READ-ONLY; The time when the disk access was created.
	TimeCreated *date.Time `json:"timeCreated,omitempty"`
	// PrivateEndpointConnections - READ-ONLY; A readonly collection of private endpoint connections created on the disk. Currently only one endpoint connection is supported.
	PrivateEndpointConnections     *[]compute.PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	PrivateEndpointsID             []string
	PrivateEndpointConnectionsID   []string
	PrivateEndpointConnectionsName []string
	PrivateEndpointConnectionsType []string
}

//// LIST FUNCTION

func listAzureComputeDiskAccesses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeDiskAccesses")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewDiskAccessesClient(subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, diskAccesss := range result.Values() {
		var PrivateEndpointID []string
		var PrivateEndpointConnectionsID []string
		var PrivateEndpointConnectionsName []string
		var PrivateEndpointConnectionsType []string
		for _, connection := range *diskAccesss.DiskAccessProperties.PrivateEndpointConnections {
			PrivateEndpointConnectionsID = append(PrivateEndpointConnectionsID, *connection.ID)
			PrivateEndpointConnectionsName = append(PrivateEndpointConnectionsName, *connection.Name)
			PrivateEndpointConnectionsType = append(PrivateEndpointConnectionsType, *connection.Type)
			PrivateEndpointID = append(PrivateEndpointID, *connection.PrivateEndpoint.ID)
		}
		d.StreamListItem(ctx, diskAccesssInfo{diskAccesss.ID, diskAccesss.Name, diskAccesss.Type, diskAccesss.Location, diskAccesss.Tags, diskAccesss.DiskAccessProperties.ProvisioningState, diskAccesss.DiskAccessProperties.TimeCreated, diskAccesss.DiskAccessProperties.PrivateEndpointConnections, PrivateEndpointID, PrivateEndpointConnectionsID, PrivateEndpointConnectionsName, PrivateEndpointConnectionsType})
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, diskAccesss := range result.Values() {
			var PrivateEndpointID []string
			var PrivateEndpointConnectionsID []string
			var PrivateEndpointConnectionsName []string
			var PrivateEndpointConnectionsType []string
			for _, connection := range *diskAccesss.DiskAccessProperties.PrivateEndpointConnections {
				PrivateEndpointConnectionsID = append(PrivateEndpointConnectionsID, *connection.ID)
				PrivateEndpointConnectionsName = append(PrivateEndpointConnectionsName, *connection.Name)
				PrivateEndpointConnectionsType = append(PrivateEndpointConnectionsType, *connection.Type)
				PrivateEndpointID = append(PrivateEndpointID, *connection.PrivateEndpoint.ID)
			}
			d.StreamListItem(ctx, diskAccesssInfo{diskAccesss.ID, diskAccesss.Name, diskAccesss.Type, diskAccesss.Location, diskAccesss.Tags, diskAccesss.DiskAccessProperties.ProvisioningState, diskAccesss.DiskAccessProperties.TimeCreated, diskAccesss.DiskAccessProperties.PrivateEndpointConnections, PrivateEndpointID, PrivateEndpointConnectionsID, PrivateEndpointConnectionsName, PrivateEndpointConnectionsType})
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAzureComputeDiskAccess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeDiskAccess")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewDiskAccessesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	diskAccesss, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if diskAccesss.ID != nil {
		var PrivateEndpointID []string
		var PrivateEndpointConnectionsID []string
		var PrivateEndpointConnectionsName []string
		var PrivateEndpointConnectionsType []string
		for _, connection := range *diskAccesss.DiskAccessProperties.PrivateEndpointConnections {
			PrivateEndpointConnectionsID = append(PrivateEndpointConnectionsID, *connection.ID)
			PrivateEndpointConnectionsName = append(PrivateEndpointConnectionsName, *connection.Name)
			PrivateEndpointConnectionsType = append(PrivateEndpointConnectionsType, *connection.Type)
			PrivateEndpointID = append(PrivateEndpointID, *connection.PrivateEndpoint.ID)
		}
		return diskAccesssInfo{diskAccesss.ID, diskAccesss.Name, diskAccesss.Type, diskAccesss.Location, diskAccesss.Tags, diskAccesss.DiskAccessProperties.ProvisioningState, diskAccesss.DiskAccessProperties.TimeCreated, diskAccesss.DiskAccessProperties.PrivateEndpointConnections, PrivateEndpointID, PrivateEndpointConnectionsID, PrivateEndpointConnectionsName, PrivateEndpointConnectionsType}, nil
	}

	return nil, nil
}
