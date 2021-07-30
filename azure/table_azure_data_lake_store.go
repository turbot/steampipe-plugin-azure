package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureDataLakeStore(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_lake_store",
		Description: "Azure Data Lake Store",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getDataLakeStore,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listDataLakeStores,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_id",
				Description: "The unique identifier associated with this Data Lake Store account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.AccountID", "DataLakeStoreAccountProperties.AccountID"),
			},
			{
				Name:        "creation_time",
				Description: "The account creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.CreationTime").Transform(convertDateToTime),
			},
			{
				Name:        "endpoint",
				Description: "The full CName endpoint for this account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning status of the Data Lake Store account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.ProvisioningState"),
			},
			{
				Name:        "state",
				Description: "The state of the Data Lake Store account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.State"),
			},
			{
				Name:        "last_modified_time",
				Description: "Managed service identity of the factory.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.LastModifiedTime").Transform(convertDateToTime),
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
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard column
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listDataLakeStores(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := account.NewAccountsClient(subscriptionID)
	accountClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := accountClient.List(context.Background(), "", nil, nil, "", "", nil)
		if err != nil {
			return nil, err
		}
		for _, account := range result.Values() {
			d.StreamListItem(ctx, account)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDataLakeStore(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataLakeStore")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := account.NewAccountsClient(subscriptionID)
	accountClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := accountClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
