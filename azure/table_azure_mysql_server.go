package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
)

//// TABLE DEFINITION

func tableAzureMySQLServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mysql_server",
		Description: "Azure MySQL Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMySQLServer,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMySQLServers,
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
				Description: "The resource type of the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.State"),
			},
			{
				Name:        "version",
				Description: "Specifies the version of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Version").Transform(transform.ToString),
			},
			{
				Name:        "user_visible_state",
				Description: "A state of a server that is visible to user. Possible values include: 'Ready', 'Dropping', 'Disabled', 'Inaccessible'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.UserVisibleState").Transform(transform.ToString),
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
				Description: "Status showing whether the server enabled infrastructure encryption. Possible values include: 'Enabled', 'Disabled'.",
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
				Description: "Indicates whether or not public network access is allowed for this server. Value is optional but if passed in, must be 'Enabled' or 'Disabled'. Possible values include: 'Enabled', 'Disabled'.",
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
				Transform:   transform.From(extractMySQLServerPrivateEndpointConnections),
			},
			{
				Name:        "server_security_alert_policy",
				Description: "Security alert policy associated with the MySQL Server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMySQLServerSecurityAlertPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_configurations",
				Description: "The server configurations(parameters) details of the server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMySQLServersConfigurations,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_keys",
				Description: "The server keys of the server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMySQLServersServerKeys,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "vnet_rules",
				Description: "Rules represented by VNET.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMySQLServerVnetRules,
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

//// LIST FUNCTION

func listMySQLServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysql.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listMySQLServers", "list", err)
		return nil, err
	}

	// Currently the API does not support pagination
	for _, server := range *result.Value {
		d.StreamListItem(ctx, server)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMySQLServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMySQLServer")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Error: mysql.ServersClient#Get: Invalid input: autorest/validation: validation failed: parameter=resourceGroupName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(resourceGroup) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysql.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getMySQLServer", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func listMySQLServersServerKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMySQLServersServerKeys")

	namespace := h.Item.(mysql.Server)
	resourceGroup := strings.Split(string(*namespace.ID), "/")[4]
	serverName := *namespace.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysql.NewServerKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, resourceGroup, serverName)
	if err != nil {
		plugin.Logger(ctx).Error("listMySQLServersServerKeys", "list", err)
		return nil, err
	}

	var mySQLServersServerKeys []map[string]interface{}

	for _, i := range op.Values() {
		mySQLServersServerKeys = append(mySQLServersServerKeys, extractMySQLServersServerKey(i))
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listMySQLServersServerKeys", "list_paging", err)
			return nil, err
		}
		for _, i := range op.Values() {
			mySQLServersServerKeys = append(mySQLServersServerKeys, extractMySQLServersServerKey(i))
		}
	}

	return mySQLServersServerKeys, nil
}

func listMySQLServerVnetRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	namespace := h.Item.(mysql.Server)
	resourceGroup := strings.Split(string(*namespace.ID), "/")[4]
	serverName := *namespace.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_server.listMySQLServerVnetRules", "connection_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysql.NewVirtualNetworkRulesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroup, serverName)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mysql_server.listMySQLServerVnetRules", "api_error", err)
		return nil, err
	}

	var vnetRules []mysql.VirtualNetworkRule

	vnetRules = append(vnetRules, op.Values()...)

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mysql_server.listMySQLServerVnetRules", "api_pagging_error", err)
			return nil, err
		}
		vnetRules = append(vnetRules, op.Values()...)
	}

	return vnetRules, nil
}

func listMySQLServersConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMySQLServersConfigurations")

	server := h.Item.(mysql.Server)
	resourceGroup := strings.Split(string(*server.ID), "/")[4]
	serverName := *server.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := mysql.NewConfigurationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroup, serverName)
	if err != nil {
		plugin.Logger(ctx).Error("listMySQLServersConfigurations", "list", err)
		return nil, err
	}

	var mySQLServersConfigurations []map[string]interface{}

	for _, i := range *op.Value {
		mySQLServersConfigurations = append(mySQLServersConfigurations, extractMySQLServersconfiguration(i))
	}

	return mySQLServersConfigurations, nil
}

func getMySQLServerSecurityAlertPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getMySQLServerSecurityAlertPolicy")

	server := h.Item.(mysql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]
	serverName := *server.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("getMySQLServerSecurityAlertPolicy")
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	

	client := mysql.NewServerSecurityAlertPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, serverName)
	if err != nil {
		return nil, err
	}

	return *op.SecurityAlertPolicyProperties, nil
}

//// TRANSFORM FUNCTION

// If we return the API response directly, the output will not provide the properties of PrivateEndpointConnections
func extractMySQLServerPrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	server := d.HydrateItem.(mysql.Server)
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

// If we return the API response directly, the output will not provide the properties of ServerKeys
func extractMySQLServersServerKey(i mysql.ServerKey) map[string]interface{} {
	mySQLServersServerKey := make(map[string]interface{})
	if i.ID != nil {
		mySQLServersServerKey["id"] = *i.ID
	}
	if i.Name != nil {
		mySQLServersServerKey["name"] = *i.Name
	}
	if i.Type != nil {
		mySQLServersServerKey["type"] = *i.Type
	}
	if i.Type != nil {
		mySQLServersServerKey["kind"] = *i.Kind
	}
	if i.ServerKeyProperties != nil {
		if i.ServerKeyProperties.ServerKeyType != nil {
			mySQLServersServerKey["serverKeyType"] = i.ServerKeyProperties.ServerKeyType
		}
		if i.ServerKeyProperties.URI != nil {
			mySQLServersServerKey["uri"] = i.ServerKeyProperties.URI
		}
		if i.ServerKeyProperties.CreationDate != nil {
			mySQLServersServerKey["creationDate"] = i.ServerKeyProperties.CreationDate
		}
	}
	return mySQLServersServerKey
}

// If we return the API response directly, the output will not provide the properties of Configurations
func extractMySQLServersconfiguration(i mysql.Configuration) map[string]interface{} {
	mySQLServersconfiguration := make(map[string]interface{})

	if i.ID != nil {
		mySQLServersconfiguration["ID"] = *i.ID
	}
	if i.Name != nil {
		mySQLServersconfiguration["Name"] = *i.Name
	}
	if i.Type != nil {
		mySQLServersconfiguration["Type"] = *i.Type
	}
	if i.ConfigurationProperties != nil {
		mySQLServersconfiguration["ConfigurationProperties"] = *i.ConfigurationProperties
	}

	return mySQLServersconfiguration
}
