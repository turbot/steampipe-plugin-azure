package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type sqlDatabaseInfo = struct {
	SQLDatabase   documentdb.SQLDatabaseGetResults
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION ////

func tableAzureCosmosDBSQLDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cosmosdb_sql_database",
		Description: "Azure Cosmos DB SQL Database",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"account_name", "name", "resource_group"}),
			ItemFromKey:       sqlDatabaseDataFromKey,
			Hydrate:           getCosmosDBSQLDatabase,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "NotFound"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCosmosDBAccounts,
			Hydrate:       listCosmosDBSQLDatabases,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the sql database",
			},
			{
				Name:        "account_name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the database account in which the database is created",
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a sql database uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.ID"),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.Type"),
			},
			{
				Name:        "autoscale_settings_max_throughput",
				Description: "Contains maximum throughput, the resource can scale up to",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Options.AutoscaleSettings.MaxThroughput"),
			},
			{
				Name:        "database_colls",
				Description: "A system generated property that specified the addressable path of the collections resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Colls"),
			},
			{
				Name:        "database_etag",
				Description: "A system generated property representing the resource etag required for optimistic concurrency control",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Etag"),
			},
			{
				Name:        "database_id",
				Description: "Name of the Cosmos DB SQL database",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.ID"),
			},
			{
				Name:        "database_rid",
				Description: "A system generated unique identifier for database",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Rid"),
			},
			{
				Name:        "database_ts",
				Description: "A system generated property that denotes the last updated timestamp of the resource",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Ts").Transform(transform.ToInt),
			},
			{
				Name:        "database_users",
				Description: "A system generated property that specifies the addressable path of the users resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Users"),
			},
			{
				Name:        "throughput",
				Description: "Contains the value of the Cosmos DB resource throughput or autoscaleSettings",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Options.Throughput"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SQLDatabase.Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SQLDatabase.ID").Transform(idToAkas),
			},
			{
				Name:        "region",
				Description: "The Azure region in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "resource_group",
				Description: "Name of the resource group the resource is created at",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_id",
				Description: "The Azure Subscription ID in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// BUILD HYDRATE INPUT ////

func sqlDatabaseDataFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	resourceGroup := quals["resource_group"].GetStringValue()
	accountName := quals["account_name"].GetStringValue()
	item := &sqlDatabaseInfo{
		Name:          &name,
		ResourceGroup: &resourceGroup,
		Account:       &accountName,
	}
	return item, nil
}

//// FETCH FUNCTIONS ////

func listCosmosDBSQLDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of cosmos db account
	account := h.Item.(databaseAccountInfo)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewSQLResourcesClient(subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.ListSQLDatabases(ctx, *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	for _, sqlDatabase := range *result.Value {
		resourceGroup := &strings.Split(string(*sqlDatabase.ID), "/")[4]
		d.StreamLeafListItem(ctx, sqlDatabaseInfo{sqlDatabase, account.Name, sqlDatabase.Name, resourceGroup, account.DatabaseAccount.Location})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getCosmosDBSQLDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sqlDatabaseData := h.Item.(*sqlDatabaseInfo)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	databaseAccountClient := documentdb.NewDatabaseAccountsClient(subscriptionID)
	databaseAccountClient.Authorizer = session.Authorizer

	op, err := databaseAccountClient.Get(ctx, *sqlDatabaseData.ResourceGroup, *sqlDatabaseData.Account)
	if err != nil {
		return nil, err
	}

	location := op.Location

	documentDBClient := documentdb.NewSQLResourcesClient(subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.GetSQLDatabase(ctx, *sqlDatabaseData.ResourceGroup, *sqlDatabaseData.Account, *sqlDatabaseData.Name)
	if err != nil {
		return nil, err
	}

	return sqlDatabaseInfo{result, sqlDatabaseData.Account, result.Name, sqlDatabaseData.ResourceGroup, location}, nil
}
