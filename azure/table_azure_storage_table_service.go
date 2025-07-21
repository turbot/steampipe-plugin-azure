package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type tableServiceInfo = struct {
	Table         storage.TableServiceProperties
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION ////

func tableAzureStorageTableService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_table_service",
		Description: "Azure Storage Table Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"storage_account_name", "resource_group"}),
			Hydrate:    getStorageTableService,
			Tags: map[string]string{
				"service": "Microsoft.Storage",
				"action":  "storageAccounts/tableServices/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageTableServices,
			Tags: map[string]string{
				"service": "Microsoft.Storage",
				"action":  "storageAccounts/tableServices/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the table service",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a table service uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Table.ID"),
			},
			{
				Name:        "storage_account_name",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Account"),
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Table.Type"),
			},
			{
				Name:        "cors_rules",
				Description: "A list of CORS rules",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Table.TableServicePropertiesProperties.Cors.CorsRules"),
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
				Transform:   transform.FromField("Table.ID").Transform(idToAkas),
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

//// FETCH FUNCTIONS ////

func listStorageTableServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	// Table is not supported for the account if storage type is FileStorage
	if account.Account.Kind == "FileStorage" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewTableServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &storageClient, d.Connection)

	result, err := storageClient.List(ctx, *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	for _, tableService := range *result.Value {
		d.StreamListItem(ctx, &tableServiceInfo{tableService, account.Name, tableService.Name, account.ResourceGroup, account.Account.Location})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getStorageTableService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageTableService")

	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	accountName := d.EqualsQuals["storage_account_name"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &storageClient, d.Connection)

	storageDetails, err := storageClient.GetProperties(ctx, resourceGroup, accountName, "")

	if err != nil {
		return nil, err
	}

	location := storageDetails.Location

	tableClient := storage.NewTableServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	tableClient.Authorizer = session.Authorizer

	op, err := tableClient.GetServiceProperties(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}

	return &tableServiceInfo{op, &accountName, op.Name, &resourceGroup, location}, nil
}
