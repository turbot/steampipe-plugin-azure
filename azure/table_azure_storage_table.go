package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
			KeyColumns:        plugin.AllColumns([]string{"storage_account_name", "resource_group", "name"}),
			Hydrate:           getStorageTable,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "OperationNotAllowedOnKind", "TableNotFound"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageTables,
		},
		Columns: []*plugin.Column{
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

			// Standard columns
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Table.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listStorageTables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewTableClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer

	result, err := storageClient.List(context.Background(), *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	for _, table := range result.Values() {
		d.StreamLeafListItem(ctx, &tableInfo{table, account.Name, table.Name, account.ResourceGroup, account.Account.Location})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getStorageTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageTable")

	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	accountName := d.KeyColumnQuals["storage_account_name"].GetStringValue()
	name := d.KeyColumnQuals["name"].GetStringValue()

	// length of the AccountName must be greater than, or equal to 3, and
	// length of the ResourceGroupName must be greater than 1, and
	// length of table name must be greater than, or equal to 3
	if len(accountName) < 3 || len(resourceGroup) < 1 || len(name) < 3 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewAccountsClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer

	storageDetails, err := storageClient.GetProperties(context.Background(), resourceGroup, accountName, "")

	if err != nil {
		return nil, err
	}

	location := storageDetails.Location

	tableClient := storage.NewTableClient(subscriptionID)
	tableClient.Authorizer = session.Authorizer

	op, err := tableClient.Get(context.Background(), resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	return &tableInfo{op, &accountName, op.Name, &resourceGroup, location}, nil
}
