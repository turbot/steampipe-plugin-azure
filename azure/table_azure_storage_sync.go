package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storagesync/mgmt/storagesync"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureStorageSync(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_sync",
		Description: "Azure Storage Sync",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAzureStorageSync,
			Tags: map[string]string{
				"service": "Microsoft.StorageSync",
				"action":  "storageSyncServices/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureStorageSyncs,
			Tags: map[string]string{
				"service": "Microsoft.StorageSync",
				"action":  "storageSyncServices/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource id for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the storage sync service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "incoming_traffic_policy",
				Description: "The incoming traffic policy of the storage sync service. Possible values include: 'AllowAllTraffic', 'AllowVirtualNetworksOnly'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.IncomingTrafficPolicy"),
			},
			{
				Name:        "last_operation_name",
				Description: "The last operation name of the storage sync service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.LastOperationName"),
			},
			{
				Name:        "last_workflow_id",
				Description: "The last workflow id of the storage sync service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.LastWorkflowID"),
			},
			{
				Name:        "storage_sync_service_status",
				Description: "The status of the storage sync service.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServiceProperties.StorageSyncServiceStatus"),
			},
			{
				Name:        "storage_sync_service_uid",
				Description: "The uid of the storage sync service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.StorageSyncServiceUID"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "List of private endpoint connection associated with the specified storage sync service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractStorageSyncPrivateEndpointConnections),
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

type StorageSyncPrivateEndpointConnections struct {
	PrivateEndpointPropertyID         interface{}
	PrivateLinkServiceConnectionState interface{}
	ProvisioningState                 interface{}
	ID                                *string
	Name                              *string
	Type                              *string
}

//// LIST FUNCTION

func listAzureStorageSyncs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := storagesync.NewServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	result, err := client.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listAzureStorageSyncs", "list", err)
		return nil, err
	}

	// The API doesn't support pagination
	for _, storage := range *result.Value {
		d.StreamListItem(ctx, storage)
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getAzureStorageSync(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureStorageSync")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := storagesync.NewServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getAzureStorageSync", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

// If we return the API response directly, the output will not provide all the properties of PrivateEndpointConnections
func extractStorageSyncPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	service := d.HydrateItem.(storagesync.Service)
	info := []StorageSyncPrivateEndpointConnections{}

	if service.ServiceProperties != nil && service.ServiceProperties.PrivateEndpointConnections != nil {
		for _, connection := range *service.ServiceProperties.PrivateEndpointConnections {
			properties := StorageSyncPrivateEndpointConnections{}
			properties.ID = connection.ID
			properties.Name = connection.Name
			properties.Type = connection.Type
			if connection.PrivateEndpointConnectionProperties != nil {
				if connection.PrivateEndpointConnectionProperties.PrivateEndpoint != nil {
					properties.PrivateEndpointPropertyID = connection.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
				}
				properties.PrivateLinkServiceConnectionState = connection.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState
				properties.ProvisioningState = connection.PrivateEndpointConnectionProperties.ProvisioningState
			}
			info = append(info, properties)
		}
	}

	return info, nil
}
