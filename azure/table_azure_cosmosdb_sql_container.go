package azure

import (
	"context"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cosmos-db/mgmt/documentdb"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureCosmosDBSQLContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cosmosdb_sql_container",
		Description: "Azure Cosmos DB SQL Container",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.AllColumns([]string{"account_name", "name", "resource_group"}),
		// 	Hydrate:    getCosmosDBSQLContainer,
		// 	IgnoreConfig: &plugin.IgnoreConfig{
		// 		ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "NotFound"}),
		// 	},
		// },
		List: &plugin.ListConfig{
			ParentHydrate: listCosmosDBAccounts,
			Hydrate:       listCosmosDBSQLContainer,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the ARM resource.",
			},
			{
				Name:        "account_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the database account.",
				Transform:   transform.FromP(extractParentPropertiesForContainer, "AccountName"),
			},
			{
				Name:        "database_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the database.",
				Transform:   transform.FromP(extractParentPropertiesForContainer, "DatabaseName"),
			},
			// {
			// 	Name:        "account_name",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "The friendly name that identifies the database account in which the database is created.",
			// 	Transform:   transform.FromField("Account"),
			// },
			{
				Name:        "id",
				Description: "The unique resource identifier of the ARM resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource",
				Description: "The resource info.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SQLContainerGetProperties.Resource"),
			},
			{
				Name:        "options",
				Description: "The resource options.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SQLContainerGetProperties.Options"),
			},
			{
				Name:        "throughput_settings",
				Description: "The resource options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCosmosDBSQLContainerThroughput,
				Transform:   transform.FromValue().Transform(transform.NullIfZeroValue),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
				Transform:   transform.FromP(extractParentPropertiesForContainer, "ResourceGroup"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCosmosDBSQLContainer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of cosmos db account
	account := h.Item.(databaseAccountInfo)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewSQLResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	// List all the database
	dbResp, err := documentDBClient.ListSQLDatabases(ctx, *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	var dbs []*documentdb.SQLDatabaseGetResults
	for _, sqlDatabase := range *dbResp.Value {
		dbs = append(dbs, &sqlDatabase)
	}

	for _, db := range dbs {
		result, err := documentDBClient.ListSQLContainers(ctx, *account.ResourceGroup, *account.Name, *db.Name)
		if err != nil {
			return nil, err
		}

		for _, sqlDatabaseContaniner := range *result.Value {
			// resourceGroup := &strings.Split(string(*sqlDatabase.ID), "/")[4]
			// data := structToMap(reflect.ValueOf(sqlDatabaseContaniner))
			d.StreamLeafListItem(ctx, sqlDatabaseContaniner)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

// func getCosmosDBSQLContainer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getCosmosDBSQLDatabase")

// 	session, err := GetNewSession(ctx, d, "MANAGEMENT")
// 	if err != nil {
// 		return nil, err
// 	}
// 	subscriptionID := session.SubscriptionID

// 	name := d.EqualsQuals["name"].GetStringValue()
// 	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
// 	accountName := d.EqualsQuals["account_name"].GetStringValue()
// 	databaseName := d.EqualsQuals["database_name"].GetStringValue()

// 	databaseAccountClient := documentdb.NewDatabaseAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
// 	databaseAccountClient.Authorizer = session.Authorizer

// 	documentDBClient := documentdb.NewSQLResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
// 	documentDBClient.Authorizer = session.Authorizer

// 	result, err := documentDBClient.GetSQLContainer(ctx, resourceGroup, accountName, databaseName, name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

func getCosmosDBSQLContainerThroughput(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of cosmos db account
	account := h.Item.(documentdb.SQLContainerGetResults)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewSQLResourcesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	p := extractContainerParentProperty(account)
	if p == nil {
		return nil, nil
	}

	result, err := documentDBClient.GetSQLContainerThroughput(ctx, p.ResourceGroup, p.AccountName, p.DatabaseName, p.Name)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		return nil, err
	}

	resultMap := structToMap(reflect.ValueOf(*result.ThroughputSettingsGetProperties))

	return resultMap, err
}

//// TRANSFORM FUNCTION

func extractParentPropertiesForContainer(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	data := d.HydrateItem.(documentdb.SQLContainerGetResults)
	p := extractContainerParentProperty(data)
	switch param {
		case "AccountName":
			return p.AccountName, nil
		case "DatabaseName":
			return p.DatabaseName, nil
		case "ResourceGroup":
			return p.ResourceGroup, nil
	}

	return nil, nil
}

//// UTILITY FUNCTION

type containerParenteInfo struct {
	ResourceGroup string
	AccountName   string
	DatabaseName  string
	Name          string
}

func extractContainerParentProperty(c documentdb.SQLContainerGetResults) *containerParenteInfo {
	if c.ID == nil {
		return nil
	}
	splitData := strings.Split(*c.ID, "/")

	// /subscriptions/d46d7416-f95f-4771-bbb5-529d4c76659c/resourceGroups/demo/providers/Microsoft.DocumentDB/databaseAccounts/test93/sqlDatabases/test63/containers/test53
	return &containerParenteInfo{
		ResourceGroup: splitData[4],
		AccountName:   splitData[8],
		DatabaseName:  splitData[10],
		Name:          splitData[12],
	}

}
