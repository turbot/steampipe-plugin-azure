package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

type databaseAccountInfo = struct {
	DatabaseAccount documentdb.DatabaseAccountGetResults
	Name            *string
	ResourceGroup   *string
}

//// TABLE DEFINITION

func tableAzureCosmosDBAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cosmosdb_account",
		Description: "Azure Cosmos DB Account",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getCosmosDBAccount,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listCosmosDBAccounts,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the database account.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a database account uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.ID"),
			},
			{
				Name:        "kind",
				Description: "Indicates the type of database account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.Kind").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.Type"),
			},
			{
				Name:        "connector_offer",
				Description: "The cassandra connector offer type for the Cosmos DB database C* account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.ConnectorOffer").Transform(transform.ToString),
			},
			{
				Name:        "consistency_policy_max_interval",
				Description: "The time amount of staleness (in seconds) tolerated, when used with the Bounded Staleness consistency level.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.ConsistencyPolicy.MaxIntervalInSeconds"),
			},
			{
				Name:        "consistency_policy_max_staleness_prefix",
				Description: "The number of stale requests tolerated, when used with the Bounded Staleness consistency level.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.ConsistencyPolicy.MaxStalenessPrefix"),
			},
			{
				Name:        "database_account_offer_type",
				Description: "The offer type for the Cosmos DB database account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.DatabaseAccountOfferType").Transform(transform.ToString),
			},
			{
				Name:        "default_consistency_level",
				Description: "The default consistency level and configuration settings of the Cosmos DB account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.ConsistencyPolicy.DefaultConsistencyLevel").Transform(transform.ToString),
			},
			{
				Name:        "disable_key_based_metadata_write_access",
				Description: "Disable write operations on metadata resources (databases, containers, throughput) via account keys.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.DisableKeyBasedMetadataWriteAccess"),
				Default:     false,
			},
			{
				Name:        "document_endpoint",
				Description: "The connection endpoint for the Cosmos DB database account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.DocumentEndpoint"),
			},
			{
				Name:        "enable_analytical_storage",
				Description: "Specifies whether to enable storage analytics, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.EnableAnalyticalStorage"),
				Default:     false,
			},
			{
				Name:        "enable_automatic_failover",
				Description: "Enables automatic failover of the write region in the rare event that the region is unavailable due to an outage.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.EnableAutomaticFailover"),
				Default:     false,
			},
			{
				Name:        "enable_cassandra_connector",
				Description: "Enables the cassandra connector on the Cosmos DB C* account.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.EnableCassandraConnector"),
				Default:     false,
			},
			{
				Name:        "enable_free_tier",
				Description: "Specifies whether free Tier is enabled for Cosmos DB database account, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.EnableFreeTier"),
				Default:     false,
			},
			{
				Name:        "enable_multiple_write_locations",
				Description: "Enables the account to write in multiple locations.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.EnableMultipleWriteLocations"),
				Default:     false,
			},
			{
				Name:        "is_virtual_network_filter_enabled",
				Description: "Specifies whether to enable/disable Virtual Network ACL rules.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.IsVirtualNetworkFilterEnabled"),
				Default:     false,
			},
			{
				Name:        "key_vault_key_uri",
				Description: "The URI of the key vault, used to encrypt the Cosmos DB database account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.KeyVaultKeyURI"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the database account resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.ProvisioningState"),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether requests from Public Network are allowed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "server_version",
				Description: "Describes the ServerVersion of an a MongoDB account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.APIProperties.ServerVersion").Transform(transform.ToString),
			},
			{
				Name:        "capabilities",
				Description: "A list of Cosmos DB capabilities for the account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.Capabilities"),
			},
			{
				Name:        "cors",
				Description: "A list of CORS policy for the Cosmos DB database account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.Cors"),
			},
			{
				Name:        "failover_policies",
				Description: "A list of regions ordered by their failover priorities.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.FailoverPolicies"),
			},
			{
				Name:        "ip_rules",
				Description: "A list of IP rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.IPRules"),
			},
			{
				Name:        "locations",
				Description: "A list of all locations that are enabled for the Cosmos DB account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.Locations"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "A list of Private Endpoint Connections configured for the Cosmos DB account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.PrivateEndpointConnections"),
			},
			{
				Name:        "read_locations",
				Description: "A list of read locations enabled for the Cosmos DB account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.ReadLocations"),
			},
			{
				Name:        "virtual_network_rules",
				Description: "A list of Virtual Network ACL rules configured for the Cosmos DB account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractCosmosDBVirtualNetworkRule),
			},
			{
				Name:        "write_locations",
				Description: "A list of write locations enabled for the Cosmos DB account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.DatabaseAccountGetProperties.WriteLocations"),
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
				Transform:   transform.FromField("DatabaseAccount.Tags"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseAccount.ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseAccount.Location").Transform(toLower),
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

func listCosmosDBAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	documentDBClient := documentdb.NewDatabaseAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	result, err := documentDBClient.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, account := range *result.Value {
		resourceGroup := &strings.Split(string(*account.ID), "/")[4]
		d.StreamListItem(ctx, databaseAccountInfo{account, account.Name, resourceGroup})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCosmosDBAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCosmosDBAccount")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	documentDBClient := documentdb.NewDatabaseAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	documentDBClient.Authorizer = session.Authorizer

	op, err := documentDBClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return databaseAccountInfo{op, op.Name, &resourceGroup}, nil
}

//// TRANSFORM FUNCTIONS

func extractCosmosDBVirtualNetworkRule(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	info := d.HydrateItem.(databaseAccountInfo)
	if info.DatabaseAccount.DatabaseAccountGetProperties != nil {
		if info.DatabaseAccount.DatabaseAccountGetProperties.VirtualNetworkRules != nil {
			return *info.DatabaseAccount.DatabaseAccountGetProperties.VirtualNetworkRules, nil
		}
	}
	return nil, nil
}
