package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2021-04-01-preview/documentdb"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type mongoCollectionInfo = struct {
	MongoCollection documentdb.MongoDBCollectionGetResults
	Account         *string
	Database        *string
	Name            *string
	ResourceGroup   *string
	Location        *string
}

//// TABLE DEFINITION

func tableAzureCosmosDBMongoCollection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cosmosdb_mongo_collection",
		Description: "Azure Cosmos DB Mongo Collection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"account_name", "name", "resource_group", "database_name"}),
			Hydrate:    getCosmosDBMongoCollection,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "NotFound"}),
			},
		},
		List: &plugin.ListConfig{
			KeyColumns:    plugin.SingleColumn("database_name"),
			ParentHydrate: listCosmosDBAccounts,
			Hydrate:       listCosmosDBMongoCollections,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the Mongo DB database.",
			},
			{
				Name:        "account_name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the database account in which the database is created.",
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "database_name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the database account in which the database is created.",
				Transform:   transform.FromField("Database"),
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a Mongo DB database uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoCollection.ID"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoCollection.Type"),
			},
			{
				Name:        "analytical_storage_ttl",
				Description: "Contains maximum throughput, the resource can scale up to.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Resource.AnalyticalStorageTTL"),
			},
			{
				Name:        "autoscale_settings_max_throughput",
				Description: "Contains maximum throughput, the resource can scale up to.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Options.AutoscaleSettings.MaxThroughput"),
			},
			{
				Name:        "collection_etag",
				Description: "A system generated property representing the resource etag required for optimistic concurrency control.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Resource.Etag"),
			},
			{
				Name:        "collection_id",
				Description: "Name of the Cosmos DB MongoDB database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Resource.ID"),
			},
			{
				Name:        "collection_rid",
				Description: "A system generated unique identifier for database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Resource.Rid"),
			},
			{
				Name:        "collection_ts",
				Description: "A system generated property that denotes the last updated timestamp of the resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Resource.Ts").Transform(transform.ToInt),
			},
			{
				Name:        "shard_key",
				Description: "A key-value pair of shard keys to be applied for the request.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Resource.ShardKey"),
			},
			{
				Name:        "indexes",
				Description: "List of index keys.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Resource.Indexes"),
			},
			{
				Name:        "throughput",
				Description: "Contains the value of the Cosmos DB resource throughput or autoscaleSettings.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MongoCollection.MongoDBCollectionGetProperties.Options.Throughput"),
			},
			{
				Name:        "throughput_settings",
				Description: "Contains the value of the Cosmos DB resource throughput or autoscaleSettings.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCosmosDBMongoCollectionThroughput,
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
				Transform:   transform.FromField("MongoCollection.Tags"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MongoCollection.ID").Transform(idToAkas),
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
				Transform:   transform.FromField("ResourceGroup").Transform(toLower),
			},
		}),
	}
}

//// LIST FUNCTION

func listCosmosDBMongoCollections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of cosmos db account
	account := h.Item.(databaseAccountInfo)
	databaseName := d.EqualsQuals["database_name"].GetStringValue()

	if databaseName == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewMongoDBResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.ListMongoDBCollections(ctx, *account.ResourceGroup, *account.Name, databaseName)
	if err != nil {
		return nil, err
	}

	for _, mongoCollection := range *result.Value {
		resourceGroup := &strings.Split(string(*mongoCollection.ID), "/")[4]
		d.StreamLeafListItem(ctx, mongoCollectionInfo{mongoCollection, account.Name, &databaseName, mongoCollection.Name, resourceGroup, mongoCollection.Location})

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCosmosDBMongoCollection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCosmosDBMongoCollection")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	accountName := d.EqualsQuals["account_name"].GetStringValue()
	databaseName := d.EqualsQuals["database_name"].GetStringValue()

	// Length of Account name must be greater than, or equal to 3
	// Error: pq: rpc error: code = Unknown desc = documentdb.DatabaseAccountsClient#Get: Invalid input: autorest/validation: validation failed: parameter=accountName
	// constraint=MinLength value="" details: value length must be greater than or equal to 3
	if len(accountName) < 3 || len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewMongoDBResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.GetMongoDBCollection(ctx, resourceGroup, accountName, databaseName, name)
	if err != nil {
		return nil, err
	}

	return mongoCollectionInfo{result, &accountName, &databaseName, result.Name, &resourceGroup, result.Location}, nil
}

func getCosmosDBMongoCollectionThroughput(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	collection := h.Item.(mongoCollectionInfo)
	databaseName := collection.Database
	resourceGroup := collection.ResourceGroup
	accountName := collection.Account
	collectionName := collection.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewMongoDBResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.GetMongoDBCollectionThroughput(ctx, *resourceGroup, *accountName, *databaseName, *collectionName)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}

	return mapThroughputSettings(result), nil
}

func mapThroughputSettings(result documentdb.ThroughputSettingsGetResults) *ThroughputSettings {
	var data ThroughputSettings

	if result.ID == nil {
		return nil
	}

	if result.ID != nil {
		data.ID = *result.ID
	}
	if result.Name != nil {
		data.Name = *result.Name
	}
	if result.Type != nil {
		data.Type = *result.Type
	}
	if result.Location != nil {
		data.Location = *result.Location
	}

	if result.Resource != nil {

		if result.Resource.Throughput != nil {
			data.Throughput = *result.Resource.Throughput
		}
		if result.Resource.AutoscaleSettings != nil {

			if result.Resource.AutoscaleSettings.MaxThroughput != nil {
				data.MaxThroughput = *result.Resource.AutoscaleSettings.MaxThroughput
			}

			if result.Resource.AutoscaleSettings.AutoUpgradePolicy.ThroughputPolicy != nil {
				data.ThroughputPolicy = documentdb.ThroughputPolicyResource{
					IsEnabled:        result.Resource.AutoscaleSettings.AutoUpgradePolicy.ThroughputPolicy.IsEnabled,
					IncrementPercent: result.Resource.AutoscaleSettings.AutoUpgradePolicy.ThroughputPolicy.IncrementPercent,
				}
			}

			if result.Resource.AutoscaleSettings.TargetMaxThroughput != nil {
				data.TargetMaxThroughput = *result.Resource.AutoscaleSettings.TargetMaxThroughput
			}
		}
		if result.Resource.MinimumThroughput != nil {
			data.MinimumThroughput = *result.Resource.MinimumThroughput
		}
		if result.Resource.OfferReplacePending != nil {
			data.OfferReplacePending = *result.Resource.OfferReplacePending
		}
		if result.Resource.Rid != nil {
			data.Rid = *result.Resource.Rid
		}
		if result.Resource.Ts != nil {
			data.Ts = *result.Resource.Ts
		}
		if result.Resource.Etag != nil {
			data.Etag = *result.Resource.Etag
		}
	}

	return &data
}
