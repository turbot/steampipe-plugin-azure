package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
)

//// TABLE DEFINITION

func tableAzureMSSQLElasticPool(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mssql_elasticpool",
		Description: "Azure Microsoft SQL Elastic Pool",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group", "server_name"}),
			Hydrate:           getMSSQLElasticPool,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listSQLServer,
			Hydrate:       listMSSQLElasticPools,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the elastic pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_name",
				Description: "The name of the parent server of the elastic pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(elasticPoolIdToServerName),
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a elastic pool uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the elastic pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the elastic pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ElasticPoolProperties.State"),
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the elastic pool.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ElasticPoolProperties.CreationDate").Transform(convertDateToTime),
			},
			{
				Name:        "database_dtu_max",
				Description: "The maximum DTU any one database can consume.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ElasticPoolProperties.DatabaseDtuMax"),
			},
			{
				Name:        "database_dtu_min",
				Description: "The minimum DTU all databases are guaranteed.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ElasticPoolProperties.DatabaseDtuMin"),
			},
			{
				Name:        "dtu",
				Description: "The total shared DTU for the database elastic pool.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ElasticPoolProperties.Dtu"),
			},
			{
				Name:        "edition",
				Description: "The edition of the elastic pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ElasticPoolProperties.Edition"),
			},
			{
				Name:        "kind",
				Description: "The kind of elastic pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_mb",
				Description: "Storage limit for the database elastic pool in MB.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ElasticPoolProperties.StorageMB"),
			},
			{
				Name:        "zone_redundant",
				Description: "Whether or not this database elastic pool is zone redundant, which means the replicas of this database will be spread across multiple availability zones.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ElasticPoolProperties.ZoneRedundant"),
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

func listMSSQLElasticPools(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := sql.NewElasticPoolsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	server := h.Item.(sql.Server)
	serverName := *server.Name
	resourceGroup := strings.Split(*server.ID, "/")[4]

	result, err := client.ListByServer(ctx, resourceGroup, serverName)
	if err != nil {
		return nil, err
	}
	for _, elasticPool := range *result.Value {
		d.StreamListItem(ctx, elasticPool)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMSSQLElasticPool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMSSQLElasticPool")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	serverName := d.KeyColumnQuals["server_name"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewElasticPoolsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
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

func elasticPoolIdToServerName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(sql.ElasticPool)
	serverName := strings.Split(string(*data.ID), "/")[8]
	return serverName, nil
}
