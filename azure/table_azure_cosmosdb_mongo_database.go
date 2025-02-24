package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cosmos-db/mgmt/documentdb"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type mongoDatabaseInfo = struct {
	MongoDatabase documentdb.MongoDBDatabaseGetResults
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION

func tableAzureCosmosDBMongoDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cosmosdb_mongo_database",
		Description: "Azure Cosmos DB Mongo Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"account_name", "name", "resource_group"}),
			Hydrate:    getCosmosDBMongoDatabase,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "NotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCosmosDBAccounts,
			Hydrate:       listCosmosDBMongoDatabases,
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
				Name:        "id",
				Description: "Contains ID to identify a Mongo DB database uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoDatabase.ID"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoDatabase.Type"),
			},
			{
				Name:        "autoscale_settings_max_throughput",
				Description: "Contains maximum throughput, the resource can scale up to.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MongoDatabase.MongoDBDatabaseGetProperties.Options.AutoscaleSettings.MaxThroughput"),
			},
			{
				Name:        "database_etag",
				Description: "A system generated property representing the resource etag required for optimistic concurrency control.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoDatabase.MongoDBDatabaseGetProperties.Resource.Etag"),
			},
			{
				Name:        "database_id",
				Description: "Name of the Cosmos DB MongoDB database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoDatabase.MongoDBDatabaseGetProperties.Resource.ID"),
			},
			{
				Name:        "database_rid",
				Description: "A system generated unique identifier for database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MongoDatabase.MongoDBDatabaseGetProperties.Resource.Rid"),
			},
			{
				Name:        "database_ts",
				Description: "A system generated property that denotes the last updated timestamp of the resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MongoDatabase.MongoDBDatabaseGetProperties.Resource.Ts").Transform(transform.ToInt),
			},
			{
				Name:        "throughput",
				Description: "Contains the value of the Cosmos DB resource throughput or autoscaleSettings.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MongoDatabase.MongoDBDatabaseGetProperties.Options.Throughput"),
			},
			{
				Name:        "throughput_settings",
				Description: "Contains the value of the Cosmos DB resource throughput or autoscaleSettings.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCosmosDBMongoThroughput,
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
				Transform:   transform.FromField("MongoDatabase.Tags"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MongoDatabase.ID").Transform(idToAkas),
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

func listCosmosDBMongoDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of cosmos db account
	account := h.Item.(databaseAccountInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewMongoDBResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &documentDBClient, d.Connection)

	result, err := documentDBClient.ListMongoDBDatabases(ctx, *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	for _, mongoDatabase := range *result.Value {
		resourceGroup := &strings.Split(string(*mongoDatabase.ID), "/")[4]
		d.StreamLeafListItem(ctx, mongoDatabaseInfo{mongoDatabase, account.Name, mongoDatabase.Name, resourceGroup, account.DatabaseAccount.Location})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCosmosDBMongoDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCosmosDBMongoDatabase")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	accountName := d.EqualsQuals["account_name"].GetStringValue()

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

	databaseAccountClient := documentdb.NewDatabaseAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	databaseAccountClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &databaseAccountClient, d.Connection)

	op, err := databaseAccountClient.Get(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}

	location := op.Location

	documentDBClient := documentdb.NewMongoDBResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.GetMongoDBDatabase(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	return mongoDatabaseInfo{result, &accountName, result.Name, &resourceGroup, location}, nil
}

type ThroughputSettings = struct {
	ID                                   string
	Name                                 string
	Type                                 string
	Location                             string
	ResourceThroughput                   int32
	ResourceMinimumThroughput            string
	ResourceOfferReplacePending          string
	ResourceRid                          string
	ResourceTs                           float64
	ResourceEtag                         string
	AutoscaleSettingsMaxThroughput       int32
	AutoscaleSettingsTargetMaxThroughput int32
	AutoscaleSettingsThroughputPolicy    documentdb.ThroughputPolicyResource
}

func getCosmosDBMongoThroughput(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(mongoDatabaseInfo)
	name := database.Name
	resourceGroup := database.ResourceGroup
	accountName := database.Account

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_cosmosdb_mongo_database.getCosmosDBMongoThroughput", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewMongoDBResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &documentDBClient, d.Connection)

	result, err := documentDBClient.GetMongoDBDatabaseThroughput(ctx, *resourceGroup, *accountName, *name)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("azure_cosmosdb_mongo_database.getCosmosDBMongoThroughput", "api_error", err)
		return nil, err
	}

	return mapThroughputSettings(result), nil
}

func mapThroughputSettings(result documentdb.ThroughputSettingsGetResults) *ThroughputSettings {
	var data ThroughputSettings

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
			data.ResourceThroughput = *result.Resource.Throughput
		}
		if result.Resource.AutoscaleSettings != nil {

			if result.Resource.AutoscaleSettings.MaxThroughput != nil {
				data.AutoscaleSettingsMaxThroughput = *result.Resource.AutoscaleSettings.MaxThroughput
			}

			if result.Resource.AutoscaleSettings.AutoUpgradePolicy != nil {
				if result.Resource.AutoscaleSettings.AutoUpgradePolicy.ThroughputPolicy != nil {
					data.AutoscaleSettingsThroughputPolicy = documentdb.ThroughputPolicyResource{
						IsEnabled:        result.Resource.AutoscaleSettings.AutoUpgradePolicy.ThroughputPolicy.IsEnabled,
						IncrementPercent: result.Resource.AutoscaleSettings.AutoUpgradePolicy.ThroughputPolicy.IncrementPercent,
					}
				}
			}

			if result.Resource.AutoscaleSettings.TargetMaxThroughput != nil {
				data.AutoscaleSettingsTargetMaxThroughput = *result.Resource.AutoscaleSettings.TargetMaxThroughput
			}
		}
		if result.Resource.MinimumThroughput != nil {
			data.ResourceMinimumThroughput = *result.Resource.MinimumThroughput
		}
		if result.Resource.OfferReplacePending != nil {
			data.ResourceOfferReplacePending = *result.Resource.OfferReplacePending
		}
		if result.Resource.Rid != nil {
			data.ResourceRid = *result.Resource.Rid
		}
		if result.Resource.Ts != nil {
			data.ResourceTs = *result.Resource.Ts
		}
		if result.Resource.Etag != nil {
			data.ResourceEtag = *result.Resource.Etag
		}
	}

	return &data
}
