package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type tableInfo = struct {
	Table         storage.Table
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION

func tableAzureStorageTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_table",
		Description: "Azure Storage Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"storage_account_name", "resource_group", "name"}),
			Hydrate:    getStorageTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "OperationNotAllowedOnKind", "TableNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageTables,
			KeyColumns:    plugin.OptionalColumns([]string{"resource_group"}),
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

//// FETCH FUNCTIONS

func listStorageTables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	// Check if the query has a resource_group filter
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	if resourceGroup != "" && resourceGroup != strings.ToLower(*account.ResourceGroup) {
		return nil, nil
	}

	// Table is not supported for the account if storage type is FileStorage or BlockBlobStorage
	if account.Account.Kind == "FileStorage" || account.Account.Kind == "BlockBlobStorage" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewTableClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	storageClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &storageClient, d.Connection)

	result, err := storageClient.List(ctx, *account.ResourceGroup, *account.Name)
	if err != nil {
		/*
		* For storage account type 'Page Blob' we are getting the kind value as 'StorageV2'.
		* Storage account type 'Page Blob' does not support table, so we are getting 'FeatureNotSupportedForAccount'/'OperationNotAllowedOnKind' error.
		* With same kind(StorageV2) of storage account, we my have different type(File Share) of storage account so we need to handle this particular error.
		 */
		if strings.Contains(err.Error(), "FeatureNotSupportedForAccount") || strings.Contains(err.Error(), "OperationNotAllowedOnKind") {
			return nil, nil
		}
		return nil, err
	}

	for _, table := range result.Values() {
		d.StreamListItem(ctx, &tableInfo{table, account.Name, table.Name, account.ResourceGroup, account.Account.Location})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, table := range result.Values() {
			d.StreamListItem(ctx, table)
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

func getStorageTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageTable")

	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	accountName := d.EqualsQuals["storage_account_name"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()

	// length of the AccountName must be greater than, or equal to 3, and
	// length of the ResourceGroupName must be greater than 1, and
	// length of table name must be greater than, or equal to 3
	if len(accountName) < 3 || len(resourceGroup) < 1 || len(name) < 3 {
		return nil, nil
	}

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

	tableClient := storage.NewTableClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	tableClient.Authorizer = session.Authorizer

	op, err := tableClient.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	return &tableInfo{op, &accountName, op.Name, &resourceGroup, location}, nil
}
