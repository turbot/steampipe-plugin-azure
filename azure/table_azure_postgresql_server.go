package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/Azure/go-autorest/autorest/date"
)

//// TABLE DEFINITION

func tableAzurePostgreSqlServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_postgresql_server",
		Description: "Azure PostgreSQL Server",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getPostgreSqlServer,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listPostgreSqlServers,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a server uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the SQL server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_visible_state",
				Description: "A state of a server that is visible to user. Possible values include: 'ServerStateReady', 'ServerStateDropping', 'ServerStateDisabled', 'ServerStateInaccessible'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.UserVisibleState").Transform(transform.ToString),
			},
			{
				Name:        "version",
				Description: "Specifies the version of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Version").Transform(transform.ToString),
			},
			{
				Name:        "location",
				Description: "The resource location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrator_login",
				Description: "Specifies the username of the administrator for this server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AdministratorLogin"),
			},
			{
				Name:        "backup_retention_days",
				Description: "Backup retention days for the server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.StorageProfile.BackupRetentionDays"),
			},
			{
				Name:        "byok_enforcement",
				Description: "Status showing whether the server data encryption is enabled with customer-managed keys.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.ByokEnforcement"),
			},
			{
				Name:        "earliest_restore_date",
				Description: "Specifies the earliest restore point creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServerProperties.EarliestRestoreDate").Transform(convertDateToTime),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.FullyQualifiedDomainName"),
			},
			{
				Name:        "geo_redundant_backup",
				Description: "Indicates whether Geo-redundant is enabled, or not for server backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.StorageProfile.GeoRedundantBackup").Transform(transform.ToString),
			},
			{
				Name:        "infrastructure_encryption",
				Description: "Status showing whether the server enabled infrastructure encryption. Possible values include: 'InfrastructureEncryptionEnabled', 'InfrastructureEncryptionDisabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.InfrastructureEncryption").Transform(transform.ToString),
			},
			{
				Name:        "master_server_id",
				Description: "The master server id of a replica server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.MasterServerID"),
			},
			{
				Name:        "minimal_tls_version",
				Description: "Enforce a minimal Tls version for the server. Possible values include: 'TLS10', 'TLS11', 'TLS12', 'TLSEnforcementDisabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.MinimalTLSVersion").Transform(transform.ToString),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether or not public network access is allowed for this server. Value is optional but if passed in, must be 'Enabled' or 'Disabled'. Possible values include: 'PublicNetworkAccessEnumEnabled', 'PublicNetworkAccessEnumDisabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "replica_capacity",
				Description: "The maximum number of replicas that a master server can have.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.ReplicaCapacity"),
			},
			{
				Name:        "replication_role",
				Description: "The replication role of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.ReplicationRole"),
			},
			{
				Name:        "sku_capacity",
				Description: "The scale up/out capacity, representing server's compute units.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_family",
				Description: "The family of hardware.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Family"),
			},
			{
				Name:        "sku_name",
				Description: "The name of the sku. For example: 'B_Gen4_1', 'GP_Gen5_8'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_size",
				Description: "The size code, to be interpreted by resource as appropriate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Size"),
			},
			{
				Name:        "sku_tier",
				Description: "The tier of the particular SKU. Possible values include: 'Basic', 'GeneralPurpose', 'MemoryOptimized'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "ssl_enforcement",
				Description: "Enable ssl enforcement or not when connect to server. Possible values include: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.SslEnforcement").Transform(transform.ToString),
			},
			{
				Name:        "storage_auto_grow",
				Description: "Indicates whether storage auto grow is enabled, or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.StorageProfile.StorageAutogrow").Transform(transform.ToString),
			},
			{
				Name:        "storage_mb",
				Description: "Indicates max storage allowed for a server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServerProperties.StorageProfile.StorageMB"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "A list of private endpoint connections on a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractPostgreSqlServerPrivateEndpointConnections),
			},
			{
				Name:        "firewall_rules",
				Description: "A list of firewall rules for a server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPostgreSQLServerFirewallRules,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_administrators",
				Description: "A list of server administrators.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPostgreSQLServerAdministrator,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_configurations",
				Description: "A list of configurations for a server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPostgreSQLServerConfigurations,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_keys",
				Description: "A list of server keys for a server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listPostgreSQLServerKeys,
				Transform:   transform.FromValue(),
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

type ServerKeyInfo struct {
	ServerKeyId           *string
	ServerKeyName         *string
	ServerKeyType         *string
	ServerKeyKind         *string
	Type                  *string
	ServerKeyUri          *string
	ServerKeyCreationDate *date.Time
}

//// LIST FUNCTION

func listPostgreSqlServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := postgresql.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, server := range *result.Value {
		d.StreamListItem(ctx, server)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getPostgreSqlServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPostgreSqlServer")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Error: postgresql.ServersClient#Get: Invalid input: autorest/validation: validation failed: parameter=resourceGroupName constraint=MinLength
	// value="" details: value length must be greater than or equal to 1
	if len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := postgresql.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func getPostgreSQLServerFirewallRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPostgreSQLServerFirewallRules")
	server := h.Item.(postgresql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := postgresql.NewFirewallRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of FirewallRuleProperties
	var firewallRules []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["ID"] = i.ID
		}
		if i.Name != nil {
			objectMap["Name"] = i.Name
		}
		if i.Type != nil {
			objectMap["Type"] = i.Type
		}
		if i.FirewallRuleProperties != nil {
			objectMap["FirewallRuleProperties"] = i.FirewallRuleProperties
		}
		firewallRules = append(firewallRules, objectMap)
	}

	return firewallRules, nil
}

func listPostgreSQLServerKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listPostgreSQLServerKeys")
	server := h.Item.(postgresql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := postgresql.NewServerKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, resourceGroupName, *server.Name)
	if err != nil {
		plugin.Logger(ctx).Error("listPostgreSQLServerKeys", "List", err)
		return nil, err
	}

	var serverKeys []ServerKeyInfo

	for _, key := range op.Values() {
		keyInfo := postgreSqlServerkeyMap(key)
		serverKeys = append(serverKeys, keyInfo)
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listPostgreSQLServerKeys", "list_paging", err)
			return nil, err
		}
		for _, key := range op.Values() {
			keyInfo := postgreSqlServerkeyMap(key)
			serverKeys = append(serverKeys, keyInfo)
		}
	}

	return serverKeys, nil
}

func getPostgreSQLServerAdministrator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPostgreSQLServerAdministrator")
	server := h.Item.(postgresql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := postgresql.NewServerAdministratorsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of ServerAdministratorProperties
	var serverAdministrators []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["ID"] = i.ID
		}
		if i.Name != nil {
			objectMap["Name"] = i.Name
		}
		if i.Type != nil {
			objectMap["Type"] = i.Type
		}
		if i.ServerAdministratorProperties != nil {
			objectMap["ServerAdministratorProperties"] = i.ServerAdministratorProperties
		}
		serverAdministrators = append(serverAdministrators, objectMap)
	}
	return serverAdministrators, nil
}

func getPostgreSQLServerConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPostgreSQLServerConfigurations")
	server := h.Item.(postgresql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := postgresql.NewConfigurationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of ConfigurationProperties
	var serverParameters []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["ID"] = i.ID
		}
		if i.Name != nil {
			objectMap["Name"] = i.Name
		}
		if i.Type != nil {
			objectMap["Type"] = i.Type
		}
		if i.ConfigurationProperties != nil {
			objectMap["ConfigurationProperties"] = i.ConfigurationProperties
		}
		serverParameters = append(serverParameters, objectMap)
	}
	return serverParameters, nil
}

func postgreSqlServerkeyMap(key postgresql.ServerKey) ServerKeyInfo {
	var serverKey ServerKeyInfo
	if key.ID != nil {
		serverKey.ServerKeyId = key.ID
	}
	if key.Name != nil {
		serverKey.ServerKeyName = key.Name
	}
	if key.Type != nil {
		serverKey.Type = key.Type
	}
	if key.Kind != nil {
		serverKey.ServerKeyKind = key.Kind
	}

	if key.ServerKeyProperties != nil {
		if key.ServerKeyProperties.CreationDate != nil {
			serverKey.ServerKeyCreationDate = key.ServerKeyProperties.CreationDate
		}
		if key.ServerKeyProperties.ServerKeyType != nil {
			serverKey.ServerKeyType = key.ServerKeyProperties.ServerKeyType
		}
		if key.ServerKeyProperties.URI != nil {
			serverKey.ServerKeyUri = key.ServerKeyProperties.URI
		}
	}

	return serverKey
}

// If we return the API response directly, the output will not provide the properties of PrivateEndpointConnections
func extractPostgreSqlServerPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	server := d.HydrateItem.(postgresql.Server)
	var properties []map[string]interface{}

	if server.ServerProperties.PrivateEndpointConnections != nil {
		for _, i := range *server.ServerProperties.PrivateEndpointConnections {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Properties != nil {
				if i.Properties.PrivateEndpoint != nil {
					objectMap["privateEndpointPropertyId"] = i.Properties.PrivateEndpoint.ID
				}
				if i.Properties.PrivateLinkServiceConnectionState != nil {
					if len(i.Properties.PrivateLinkServiceConnectionState.ActionsRequired) > 0 {
						objectMap["privateLinkServiceConnectionStateActionsRequired"] = i.Properties.PrivateLinkServiceConnectionState.ActionsRequired
					}
					if len(i.Properties.PrivateLinkServiceConnectionState.Status) > 0 {
						objectMap["privateLinkServiceConnectionStateStatus"] = i.Properties.PrivateLinkServiceConnectionState.Status
					}
					if i.Properties.PrivateLinkServiceConnectionState.Description != nil {
						objectMap["privateLinkServiceConnectionStateDescription"] = i.Properties.PrivateLinkServiceConnectionState.Description
					}
				}
				if len(i.Properties.ProvisioningState) > 0 {
					objectMap["provisioningState"] = i.Properties.ProvisioningState
				}
			}
			properties = append(properties, objectMap)
		}
	}

	return properties, nil
}
