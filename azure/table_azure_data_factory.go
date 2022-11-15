package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/datafactory/mgmt/datafactory"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAzureDataFactory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_factory",
		Description: "Azure Data Factory",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getDataFactory,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listDataFactories,
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
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "Version of the factory.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FactoryProperties.Version"),
			},
			{
				Name:        "create_time",
				Description: "Specifies the time, the factory was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("FactoryProperties.CreateTime").Transform(convertDateToTime),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "provisioning_state",
				Description: "Factory provisioning state, example Succeeded.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FactoryProperties.ProvisioningState"),
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public network access is allowed for the data factory.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FactoryProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "additional_properties",
				Description: "Unmatched properties from the message are deserialized this collection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "Managed service identity of the factory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "encryption",
				Description: "Properties to enable Customer Managed Key for the factory.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FactoryProperties.EncryptionConfiguration"),
			},
			{
				Name:        "repo_configuration",
				Description: "Git repo information of the factory.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FactoryProperties.RepoConfiguration"),
			},
			{
				Name:        "global_parameters",
				Description: "List of parameters for factory.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FactoryProperties.GlobalParameters"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "List of private endpoint connections for data factory.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listDataFactoryPrivateEndpointConnections,
				Transform:   transform.FromValue(),
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

func listDataFactories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	factoryClient := datafactory.NewFactoriesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	factoryClient.Authorizer = session.Authorizer

	result, err := factoryClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, factory := range result.Values() {
		d.StreamListItem(ctx, factory)
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
		for _, factory := range result.Values() {
			d.StreamListItem(ctx, factory)
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

func getDataFactory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataFactory")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	factoryClient := datafactory.NewFactoriesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	factoryClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := factoryClient.Get(ctx, resourceGroup, name, "*")
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listDataFactoryPrivateEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDataFactoryPrivateEndpointConnections")
	factory := h.Item.(datafactory.Factory)
	factoryName := factory.Name
	resourceGroup := strings.Split(*factory.ID, "/")[4]

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("listDataFactoryPrivateEndpointConnections", "connection", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	connClient := datafactory.NewPrivateEndPointConnectionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	connClient.Authorizer = session.Authorizer

	op, err := connClient.ListByFactory(ctx, resourceGroup, *factoryName)
	if err != nil {
		plugin.Logger(ctx).Error("listDataFactoryPrivateEndpointConnections", "ListByFactory", err)
		return nil, err
	}

	var connections []PrivateConnection
	var connection PrivateConnection

	for _, conn := range op.Values() {
		connection = factoryPrivateEndpointConnectionMap(conn)
		connections = append(connections, connection)
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listDataFactoryPrivateEndpointConnections", "ListByFactory_pagination", err)
			return nil, err
		}
		for _, conn := range op.Values() {
			connection = factoryPrivateEndpointConnectionMap(conn)
			connections = append(connections, connection)
		}
	}

	return connections, nil
}

// If we return the API response directly, the output will not give
// all the properties of PrivateEndpointConnection
func factoryPrivateEndpointConnectionMap(conn datafactory.PrivateEndpointConnectionResource) PrivateConnection {
	var connection PrivateConnection
	if conn.ID != nil {
		connection.PrivateEndpointConnectionId = conn.ID
	}
	if conn.Name != nil {
		connection.PrivateEndpointConnectionName = conn.Name
	}
	if conn.Type != nil {
		connection.PrivateEndpointConnectionType = conn.Type
	}
	if conn.Etag != nil {
		connection.Etag = conn.Etag
	}
	if conn.Properties != nil {
		if conn.Properties.PrivateEndpoint != nil {
			if conn.Properties.PrivateEndpoint.ID != nil {
				connection.PrivateEndpointId = conn.Properties.PrivateEndpoint.ID
			}
		}
		if conn.Properties.PrivateLinkServiceConnectionState != nil {
			if conn.Properties.PrivateLinkServiceConnectionState.ActionsRequired != nil {
				connection.PrivateLinkServiceConnectionStateActionsRequired = conn.Properties.PrivateLinkServiceConnectionState.ActionsRequired
			}
			if conn.Properties.PrivateLinkServiceConnectionState.Status != nil {
				connection.PrivateLinkServiceConnectionStateStatus = conn.Properties.PrivateLinkServiceConnectionState.Status
			}
			if conn.Properties.PrivateLinkServiceConnectionState.Description != nil {
				connection.PrivateLinkServiceConnectionStateDescription = conn.Properties.PrivateLinkServiceConnectionState.Description
			}
		}
		if conn.Properties.ProvisioningState != nil {
			connection.ProvisioningState = conn.Properties.ProvisioningState
		}
	}

	return connection
}

type PrivateConnection struct {
	ProvisioningState                                *string
	PrivateEndpointConnectionId                      *string
	PrivateEndpointId                                *string
	PrivateLinkServiceConnectionStateStatus          *string
	PrivateLinkServiceConnectionStateDescription     *string
	PrivateLinkServiceConnectionStateActionsRequired *string
	PrivateEndpointConnectionName                    *string
	PrivateEndpointConnectionType                    *string
	Etag                                             *string
}
