package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cosmos-db/mgmt/documentdb"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureCosmosDBRestorableDatabaseAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cosmosdb_restorable_database_account",
		Description: "Azure Cosmos DB Restorable Database Account",
		List: &plugin.ListConfig{
			Hydrate: listCosmosDBRestorableDatabaseAccounts,
			Tags: map[string]string{
				"service": "Microsoft.DocumentDB",
				"action":  "restorableDatabaseAccounts/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the restorable database account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a restorable database account uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "account_name",
				Description: "The name of the global database account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RestorableDatabaseAccountProperties.AccountName"),
			},
			{
				Name:        "api_type",
				Description: "The API type of the restorable database account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RestorableDatabaseAccountProperties.APIType"),
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the restorable database account.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("RestorableDatabaseAccountProperties.CreationTime").Transform(convertDateToTime),
			},
			{
				Name:        "deletion_time",
				Description: "The time at which the restorable database account has been deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("RestorableDatabaseAccountProperties.DeletionTime").Transform(convertDateToTime),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "restorable_locations",
				Description: "List of regions where the database account can be restored from.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(mapRestorableLocations),
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listCosmosDBRestorableDatabaseAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_cosmosdb_restorable_database_account.listCosmosDBRestorableDatabaseAccounts", "session_error", err)
		return nil, err
	}

	documentDBClient := documentdb.NewRestorableDatabaseAccountsClientWithBaseURI(session.ResourceManagerEndpoint, session.SubscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &documentDBClient, d.Connection)

	result, err := documentDBClient.List(ctx)
	if err != nil {
		logger.Error("azure_cosmosdb_restorable_database_account.listCosmosDBRestorableDatabaseAccounts", "api_error", err)
		return nil, err
	}

	for _, account := range *result.Value {
		d.StreamListItem(ctx, account)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

type RestorableLocationResource struct {
	LocationName                      string
	RegionalDatabaseAccountInstanceID string
	CreationTime                      date.Time
	DeletionTime                      date.Time
}

func mapRestorableLocations(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(documentdb.RestorableDatabaseAccountGetResult)

	restorableLocations := *data.RestorableDatabaseAccountProperties.RestorableLocations

	if len(restorableLocations) < 1 {
		return nil, nil
	}

	var restorableLocationsArr []RestorableLocationResource
	for _, location := range restorableLocations {

		var locationResource RestorableLocationResource

		if location.LocationName != nil {
			locationResource.LocationName = string(*location.LocationName)
		}
		if location.RegionalDatabaseAccountInstanceID != nil {
			locationResource.RegionalDatabaseAccountInstanceID = string(*location.RegionalDatabaseAccountInstanceID)
		}
		if location.CreationTime != nil {
			locationResource.CreationTime = date.Time(*location.CreationTime)
		}
		if location.DeletionTime != nil {
			locationResource.DeletionTime = date.Time(*location.DeletionTime)
		}

		restorableLocationsArr = append(restorableLocationsArr, locationResource)
	}

	return restorableLocationsArr, nil
}
