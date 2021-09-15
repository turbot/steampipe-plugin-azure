package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
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
				Transform:   transform.FromField("DiskAccessProperties.ProvisioningState"),
			},
			{
				Name:        "time_created",
				Description: "The time when the disk access was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DiskAccessProperties.TimeCreated").Transform(convertDateToTime),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The private endpoint connections details.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointConnections"),
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
	compute.DiskAccess
	// PrivateEndpointConnections - READ-ONLY; A readonly collection of private endpoint connections created on the disk.
	PrivateEndpointConnections []PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
}

type PrivateEndpointConnection struct {
	// ID - READ-ONLY; private endpoint connection Id
	ID string
	// Name - READ-ONLY; private endpoint connection name
	Name string
	// Type - READ-ONLY; private endpoint connection type
	Type string
	// PrivateEndpointID - The Id of private end point.
	PrivateEndpointID string
	// ProvisioningState - The provisioning state of the private endpoint connection resource. Possible values include: 'PrivateEndpointConnectionProvisioningStateSucceeded', 'PrivateEndpointConnectionProvisioningStateCreating', 'PrivateEndpointConnectionProvisioningStateDeleting', 'PrivateEndpointConnectionProvisioningStateFailed'
	ProvisioningState                                string
	PrivateLinkServiceConnectionStateStatus          string
	PrivateLinkServiceConnectionStateDescription     string
	PrivateLinkServiceConnectionStateActionsRequired string
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
		plugin.Logger(ctx).Error("listAzureComputeDiskAccesses", "list_err", err)
		return nil, err
	}

	// If we return the API response directly, the output will not provide
	// all the properties of diskAccesss
	for _, diskAccesss := range result.Values() {
		var PrivateEndpointConnections []PrivateEndpointConnection
		if diskAccesss.DiskAccessProperties.PrivateEndpointConnections != nil {
			for _, connection := range *diskAccesss.DiskAccessProperties.PrivateEndpointConnections {
				var PrivateConnection PrivateEndpointConnection
				if connection.ID != nil {
					PrivateConnection.ID = *connection.ID
				}
				if connection.Name != nil {
					PrivateConnection.Name = *connection.Name
				}
				if connection.Type != nil {
					PrivateConnection.Type = *connection.Type
				}
				if connection.PrivateEndpointConnectionProperties != nil {
					if connection.PrivateEndpoint != nil {
						PrivateConnection.PrivateEndpointID = *connection.PrivateEndpoint.ID
					}
					if connection.PrivateLinkServiceConnectionState != nil {
						if connection.PrivateLinkServiceConnectionState.ActionsRequired != nil {
							PrivateConnection.PrivateLinkServiceConnectionStateActionsRequired = *connection.PrivateLinkServiceConnectionState.ActionsRequired
						}
						if connection.PrivateLinkServiceConnectionState.Description != nil {
							PrivateConnection.PrivateLinkServiceConnectionStateDescription = *connection.PrivateLinkServiceConnectionState.Description
						}
						if connection.PrivateLinkServiceConnectionState.Status != "" {
							PrivateConnection.PrivateLinkServiceConnectionStateStatus = string(connection.PrivateLinkServiceConnectionState.Status)
						}
					}
					if connection.ProvisioningState != "" {
						PrivateConnection.ProvisioningState = string(connection.ProvisioningState)
					}
				}

				PrivateEndpointConnections = append(PrivateEndpointConnections, PrivateConnection)
			}
		}
		d.StreamListItem(ctx, diskAccesssInfo{diskAccesss, PrivateEndpointConnections})
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listAzureComputeDiskAccesses", "list_err", err)
			return nil, err
		}
		for _, diskAccesss := range result.Values() {
			var PrivateEndpointConnections []PrivateEndpointConnection
			if diskAccesss.DiskAccessProperties.PrivateEndpointConnections != nil {
				for _, connection := range *diskAccesss.DiskAccessProperties.PrivateEndpointConnections {
					var PrivateConnection PrivateEndpointConnection
					if connection.ID != nil {
						PrivateConnection.ID = *connection.ID
					}
					if connection.Name != nil {
						PrivateConnection.Name = *connection.Name
					}
					if connection.Type != nil {
						PrivateConnection.Type = *connection.Type
					}
					if connection.PrivateEndpointConnectionProperties != nil {
						if connection.PrivateEndpoint != nil {
							PrivateConnection.PrivateEndpointID = *connection.PrivateEndpoint.ID
						}
						if connection.PrivateLinkServiceConnectionState != nil {
							if connection.PrivateLinkServiceConnectionState.ActionsRequired != nil {
								PrivateConnection.PrivateLinkServiceConnectionStateActionsRequired = *connection.PrivateLinkServiceConnectionState.ActionsRequired
							}
							if connection.PrivateLinkServiceConnectionState.Description != nil {
								PrivateConnection.PrivateLinkServiceConnectionStateDescription = *connection.PrivateLinkServiceConnectionState.Description
							}
							if connection.PrivateLinkServiceConnectionState.Status != "" {
								PrivateConnection.PrivateLinkServiceConnectionStateStatus = string(connection.PrivateLinkServiceConnectionState.Status)
							}
						}
						if connection.ProvisioningState != "" {
							PrivateConnection.ProvisioningState = string(connection.ProvisioningState)
						}
					}

					PrivateEndpointConnections = append(PrivateEndpointConnections, PrivateConnection)
				}
			}
			d.StreamListItem(ctx, diskAccesssInfo{diskAccesss, PrivateEndpointConnections})
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAzureComputeDiskAccess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeDiskAccess")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewDiskAccessesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	diskAccesss, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("listAzureComputeDiskAccesses", "get_err", err)
		return nil, err
	}

	// If we return the API response directly, the output will not provide
	// all the properties of diskAccesss
	var PrivateEndpointConnections []PrivateEndpointConnection
	if diskAccesss.DiskAccessProperties.PrivateEndpointConnections != nil {
		for _, connection := range *diskAccesss.DiskAccessProperties.PrivateEndpointConnections {
			var PrivateConnection PrivateEndpointConnection
			if connection.ID != nil {
				PrivateConnection.ID = *connection.ID
			}
			if connection.Name != nil {
				PrivateConnection.Name = *connection.Name
			}
			if connection.Type != nil {
				PrivateConnection.Type = *connection.Type
			}
			if connection.PrivateEndpointConnectionProperties != nil {
				if connection.PrivateEndpoint != nil {
					PrivateConnection.PrivateEndpointID = *connection.PrivateEndpoint.ID
				}
				if connection.PrivateLinkServiceConnectionState != nil {
					if connection.PrivateLinkServiceConnectionState.ActionsRequired != nil {
						PrivateConnection.PrivateLinkServiceConnectionStateActionsRequired = *connection.PrivateLinkServiceConnectionState.ActionsRequired
					}
					if connection.PrivateLinkServiceConnectionState.Description != nil {
						PrivateConnection.PrivateLinkServiceConnectionStateDescription = *connection.PrivateLinkServiceConnectionState.Description
					}
					if connection.PrivateLinkServiceConnectionState.Status != "" {
						PrivateConnection.PrivateLinkServiceConnectionStateStatus = string(connection.PrivateLinkServiceConnectionState.Status)
					}
				}
				if connection.ProvisioningState != "" {
					PrivateConnection.ProvisioningState = string(connection.ProvisioningState)
				}
			}

			PrivateEndpointConnections = append(PrivateEndpointConnections, PrivateConnection)
		}
	}

	return diskAccesssInfo{diskAccesss, PrivateEndpointConnections}, nil
}
