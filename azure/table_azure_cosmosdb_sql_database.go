package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

type sqlDatabaseInfo = struct {
	SQLDatabase   documentdb.SQLDatabaseGetResults
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION

func tableAzureCosmosDBSQLDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cosmosdb_sql_database",
		Description: "Azure Cosmos DB SQL Database",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"account_name", "name", "resource_group"}),
			Hydrate:           getCosmosDBSQLDatabase,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "NotFound"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCosmosDBAccounts,
			Hydrate:       listCosmosDBSQLDatabases,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the sql database.",
			},
			{
				Name:        "account_name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the database account in which the database is created.",
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a sql database uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.ID"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.Type"),
			},
			{
				Name:        "autoscale_settings_max_throughput",
				Description: "Contains maximum throughput, the resource can scale up to.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Options.AutoscaleSettings.MaxThroughput"),
			},
			{
				Name:        "database_colls",
				Description: "A system generated property that specified the addressable path of the collections resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Colls"),
			},
			{
				Name:        "database_etag",
				Description: "A system generated property representing the resource etag required for optimistic concurrency control.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Etag"),
			},
			{
				Name:        "database_id",
				Description: "Name of the Cosmos DB SQL database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.ID"),
			},
			{
				Name:        "database_rid",
				Description: "A system generated unique identifier for database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Resource.Rid"),
			},
			{
				Name:        "database_ts",
				Description: "A system generated property that denotes the last updated timestamp of the resource.",
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
				Description: "Contains the value of the Cosmos DB resource throughput or autoscaleSettings.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SQLDatabase.SQLDatabaseGetProperties.Options.Throughput"),
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
				Transform:   transform.FromField("SQLDatabase.Tags"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SQLDatabase.ID").Transform(idToAkas),
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

func listCosmosDBSQLDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of cosmos db account
	account := h.Item.(databaseAccountInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewSQLResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.ListSQLDatabases(ctx, *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	for _, sqlDatabase := range *result.Value {
		resourceGroup := &strings.Split(string(*sqlDatabase.ID), "/")[4]
		d.StreamLeafListItem(ctx, sqlDatabaseInfo{sqlDatabase, account.Name, sqlDatabase.Name, resourceGroup, account.DatabaseAccount.Location})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCosmosDBSQLDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCosmosDBSQLDatabase")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	accountName := d.KeyColumnQuals["account_name"].GetStringValue()

	databaseAccountClient := documentdb.NewDatabaseAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	databaseAccountClient.Authorizer = session.Authorizer

	op, err := databaseAccountClient.Get(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}

	location := op.Location

	documentDBClient := documentdb.NewSQLResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.GetSQLDatabase(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	return sqlDatabaseInfo{result, &accountName, result.Name, &resourceGroup, location}, nil
}
