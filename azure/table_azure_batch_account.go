package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2020-09-01/batch"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureBatchAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_batch_account",
		Description: "Azure Batch Account",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getBatchAccount,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listBatchAccounts,
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
				Name:        "AccountEndpoint",
				Description: "The account endpoint used to interact with the Batch service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.AccountEndpoint"),
			},
			{
				Name:        "create_time",
				Description: "Specifies the time, the factory was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("FactoryProperties.CreateTime").Transform(convertDateToTime),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "provisioning_state",
				Description: "Factory provisioning state, example Succeeded.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FactoryProperties.ProvisioningState"),
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public network access is allowed for the data factory.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FactoryProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "additional_properties",
				Description: "Unmatched properties from the message are deserialized this collection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "Managed service identity of the factory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "encryption",
				Description: "Properties to enable Customer Managed Key for the factory.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FactoryProperties.EncryptionConfiguration"),
			},
			{
				Name:        "repo_configuration",
				Description: "Git repo information of the factory.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FactoryProperties.RepoConfiguration"),
			},
			{
				Name:        "global_parameters",
				Description: "List of parameters for factory.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FactoryProperties.GlobalParameters"),
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

func listBatchAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	batchAccountClient := batch.NewAccountClient(subscriptionID)
	batchAccountClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := batchAccountClient.List(context.Background())
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

func getBatchAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBatchAccount")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	batchAccountClient := batch.NewAccountClient(subscriptionID)
	batchAccountClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := batchAccountClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
