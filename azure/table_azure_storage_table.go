package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type tableInfo = struct {
	Table         storage.TableServiceProperties
	Account       *string
	Name          *string
	ResourceGroup *string
	Location      *string
}

//// TABLE DEFINITION ////

func tableAzureStorageTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_table",
		Description: "Azure Storage Table",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"storage_account_name", "resource_group"}),
			ItemFromKey:       tableDataFromKey,
			Hydrate:           getStorageTable,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageTables,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the table",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a table uniquely",
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

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Table.ID").Transform(idToAkas),
			},
			{
				Name:        "region",
				Description: "The Azure region in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "resource_group",
				Description: "Name of the resource group, the table is created at",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_id",
				Description: "The Azure Subscription ID in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Table.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// BUILD HYDRATE INPUT ////

func tableDataFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	resourceGroup := quals["resource_group"].GetStringValue()
	accountName := quals["storage_account_name"].GetStringValue()
	item := &tableInfo{
		Account:       &accountName,
		ResourceGroup: &resourceGroup,
	}
	return item, nil
}

//// FETCH FUNCTIONS ////

func listStorageTables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewTableServicesClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer

	result, err := storageClient.List(context.Background(), *account.ResourceGroup, *account.Name)
	if err != nil {
		return nil, err
	}

	for _, table := range *result.Value {
		d.StreamLeafListItem(ctx, &tableInfo{table, account.Name, table.Name, account.ResourceGroup, account.Account.Location})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getStorageTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tableData := h.Item.(*tableInfo)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	storageClient := storage.NewAccountsClient(subscriptionID)
	storageClient.Authorizer = session.Authorizer

	storageDetails, err := storageClient.GetProperties(context.Background(), *tableData.ResourceGroup, *tableData.Account, "")

	if err != nil {
		return nil, err
	}

	location := storageDetails.Location

	tableClient := storage.NewTableServicesClient(subscriptionID)
	tableClient.Authorizer = session.Authorizer

	op, err := tableClient.GetServiceProperties(context.Background(), *tableData.ResourceGroup, *tableData.Account)
	if err != nil {
		return nil, err
	}

	return &tableInfo{op, tableData.Account, op.Name, tableData.ResourceGroup, location}, nil
}
