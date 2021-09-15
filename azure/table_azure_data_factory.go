package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureDataFactory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_factory",
		Description: "Azure Data Factory",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getDataFactory,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listDataFactories,
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
				Description: "List of private connections for factory.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataFactoryPrivateConnections,
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

func listDataFactories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	factoryClient := datafactory.NewFactoriesClient(subscriptionID)
	factoryClient.Authorizer = session.Authorizer

	result, err := factoryClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, factory := range result.Values() {
		d.StreamListItem(ctx, factory)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, factory := range result.Values() {
			d.StreamListItem(ctx, factory)
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

	factoryClient := datafactory.NewFactoriesClient(subscriptionID)
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

func getDataFactoryPrivateConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataFactoryPrivateConnections")
	factory := h.Item.(datafactory.Factory)
	factoryName := factory.Name
	resourceGroup := strings.Split(*factory.ID, "/")[4]

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("getDataFactoryPrivateConnections", "connection", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	connClient := datafactory.NewPrivateEndPointConnectionsClient(subscriptionID)
	connClient.Authorizer = session.Authorizer

	op, err := connClient.ListByFactory(ctx, resourceGroup, *factoryName)
	if err != nil {
		return nil, err
	}

	var connections []PrivateConnection
	for _, connection := range op.Values() {
		plugin.Logger(ctx).Trace("Private Connections =>", connection)
		connections = append(connections, PrivateConnection{
			Properties: connection.Properties,
			Id:         connection.ID,
			Name:       connection.Name,
			Type:       connection.Type,
			Etag:       connection.Etag,
		})
	}

	return connections, nil
}

type PrivateConnection struct {
	Properties interface{}
	Id         *string
	Name       *string
	Type       *string
	Etag       *string
}
