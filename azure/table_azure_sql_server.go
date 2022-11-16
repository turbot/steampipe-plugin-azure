package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	sql "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
)

//// TABLE DEFINITION

func tableAzureSQLServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_sql_server",
		Description: "Azure SQL Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getSQLServer,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSQLServer,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the SQL server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a SQL server uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the SQL server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.State"),
			},
			{
				Name:        "kind",
				Description: "The Kind of sql server.",
				Type:        proto.ColumnType_STRING,
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
				Transform:   transform.FromField("Properties.AdministratorLogin"),
			},
			{
				Name:        "administrator_login_password",
				Description: "The administrator login password.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AdministratorLoginPassword"),
			},
			{
				Name:        "minimal_tls_version",
				Description: "Minimal TLS version. Allowed values: '1.0', '1.1', '1.2'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinimalTLSVersion"),
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public endpoint access is allowed for this server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PublicNetworkAccess"),
			},
			{
				Name:        "version",
				Description: "The version of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Version"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.FullyQualifiedDomainName"),
			},
			{
				Name:        "server_audit_policy",
				Description: "Specifies the audit policy configuration for server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLServerAuditPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_security_alert_policy",
				Description: "Specifies the security alert policy configuration for server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLServerSecurityAlertPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_azure_ad_administrator",
				Description: "Specifies the active directory administrator.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLServerAzureADAdministrator,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_vulnerability_assessment",
				Description: "Specifies the server's vulnerability assessment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLServerVulnerabilityAssessment,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "firewall_rules",
				Description: "A list of firewall rules fro this server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSQLServerFirewallRules,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption_protector",
				Description: "The server encryption protector.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLServerEncryptionProtector,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The private endpoint connections of the sql server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSQLServerPrivateEndpointConnections,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "Specifies the set of tags attached to the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "virtual_network_rules",
				Description: "A list of virtual network rules for this server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSQLServerVirtualNetworkRules,
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

type PrivateConnectionInfo struct {
	PrivateEndpointConnectionId                      string
	PrivateEndpointId                                string
	PrivateEndpointConnectionName                    string
	PrivateEndpointConnectionType                    string
	PrivateLinkServiceConnectionStateStatus          string
	PrivateLinkServiceConnectionStateDescription     string
	PrivateLinkServiceConnectionStateActionsRequired string
	ProvisioningState                                string
}

//// LIST FUNCTION

func listSQLServer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServer", "connection error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServer", "session error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewServersClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServer", "client error", err)
	}

	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.listSQLServer", "api error", err)
		}
		for _, server := range nextResult.Value {
			d.StreamListItem(ctx, server)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSQLServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServer", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewServersClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServer", "client error", err)
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServer", "api error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func getSQLServerAuditPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerAuditPolicy", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewServerBlobAuditingPoliciesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerAuditPolicy", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of ServerBlobAuditingPolicyProperties
	var auditPolicies []map[string]interface{}
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.getSQLServerAuditPolicy", "api error", err)
		}
		for _, policy := range nextResult.Value {
			objectMap := make(map[string]interface{})
			if policy.ID != nil {
				objectMap["id"] = policy.ID
			}
			if policy.Name != nil {
				objectMap["name"] = policy.Name
			}
			if policy.Type != nil {
				objectMap["type"] = policy.Type
			}
			if policy.Properties != nil {
				objectMap["properties"] = policy.Properties
			}
			auditPolicies = append(auditPolicies, objectMap)
		}
	}
	return auditPolicies, nil
}

func listSQLServerPrivateEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServerPrivateEndpointConnections", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewPrivateEndpointConnectionsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServerPrivateEndpointConnections", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}

	var privateEndpointConnections []PrivateConnectionInfo
	var connection PrivateConnectionInfo

	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.getSQLServerAuditPolicy", "api error", err)
		}
		for _, endpoint := range nextResult.Value {
			connection = privateEndpointConnectionMap(endpoint)
			privateEndpointConnections = append(privateEndpointConnections, connection)
		}
	}

	return privateEndpointConnections, nil
}

func getSQLServerSecurityAlertPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerSecurityAlertPolicy", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewServerSecurityAlertPoliciesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerSecurityAlertPolicy", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}

	var securityAlertPolicies []map[string]interface{}

	// If we return the API response directly, the output only gives
	// the contents of SecurityAlertPolicyProperties
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.getSQLServerSecurityAlertPolicy", "api error", err)
		}
		for _, alertPolicy := range nextResult.Value {
			objectMap := make(map[string]interface{})
			if alertPolicy.ID != nil {
				objectMap["id"] = alertPolicy.ID
			}
			if alertPolicy.Name != nil {
				objectMap["name"] = alertPolicy.Name
			}
			if alertPolicy.Type != nil {
				objectMap["type"] = alertPolicy.Type
			}
			if alertPolicy.Properties != nil {
				objectMap["properties"] = alertPolicy.Properties
			}
			securityAlertPolicies = append(securityAlertPolicies, objectMap)
		}
	}

	return securityAlertPolicies, nil
}

func getSQLServerAzureADAdministrator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerAzureADAdministrator", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewServerAzureADAdministratorsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerAzureADAdministrator", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}

	var serverAdministrators []map[string]interface{}
	// If we return the API response directly, the output only gives
	// the contents of ServerAdministratorProperties
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.getSQLServerAzureADAdministrator", "api error", err)
		}
		for _, i := range nextResult.Value {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Type != nil {
				objectMap["type"] = i.Type
			}
			if i.Properties != nil {
				objectMap["properties"] = i.Properties
			}
			serverAdministrators = append(serverAdministrators, objectMap)
		}
	}

	return serverAdministrators, nil
}

func getSQLServerEncryptionProtector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerEncryptionProtector", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewEncryptionProtectorsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerEncryptionProtector", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}

	var encryptionProtectors []map[string]interface{}

	// If we return the API response directly, the output only gives
	// the contents of ServerAdministratorProperties
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.getSQLServerEncryptionProtector", "api error", err)
		}
		for _, i := range nextResult.Value {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Type != nil {
				objectMap["type"] = i.Type
			}
			if i.Location != nil {
				objectMap["location"] = i.Location
			}
			if i.Kind != nil {
				objectMap["kind"] = i.Kind
			}
			if i.Properties != nil {

				if i.Properties.Subregion != nil {
					objectMap["subregion"] = i.Properties.Subregion
				}
				if i.Properties.ServerKeyName != nil {
					objectMap["serverKeyName"] = i.Properties.ServerKeyName
				}
				if i.Properties.ServerKeyType != nil {
					objectMap["serverKeyType"] = i.Properties.ServerKeyType
				}
				if i.Properties.URI != nil {
					objectMap["uri"] = i.Properties.URI
				}
				if i.Properties.Thumbprint != nil {
					objectMap["thumbprint"] = i.Properties.Thumbprint
				}
			}
			encryptionProtectors = append(encryptionProtectors, objectMap)
		}
	}

	return encryptionProtectors, nil
}

func getSQLServerVulnerabilityAssessment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerVulnerabilityAssessment", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewServerVulnerabilityAssessmentsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.getSQLServerVulnerabilityAssessment", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}

	var vulnerabilityAssessments []map[string]interface{}

	// If we return the API response directly, the output only gives
	// the contents of ServerVulnerabilityAssessmentProperties
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.getSQLServerVulnerabilityAssessment", "api error", err)
		}
		for _, i := range nextResult.Value {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Type != nil {
				objectMap["type"] = i.Type
			}
			if i.Properties != nil {
				objectMap["properties"] = i.Properties
			}
			vulnerabilityAssessments = append(vulnerabilityAssessments, objectMap)
		}
	}

	return vulnerabilityAssessments, nil
}

func listSQLServerFirewallRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServerFirewallRules", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewFirewallRulesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServerFirewallRules", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}
	var firewallRules []map[string]interface{}

	// If we return the API response directly, the output only gives
	// the contents of FirewallRuleProperties
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.listSQLServerFirewallRules", "api error", err)
		}
		for _, i := range nextResult.Value {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Type != nil {
				objectMap["type"] = i.Type
			}
			if i.Properties != nil {
				objectMap["properties"] = i.Properties
			}
			firewallRules = append(firewallRules, objectMap)
		}
	}

	return firewallRules, nil
}

func listSQLServerVirtualNetworkRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServerVirtualNetworkRules", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewVirtualNetworkRulesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_server.listSQLServerVirtualNetworkRules", "client error", err)
	}

	op := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	if err != nil {
		return nil, err
	}
	var NetworkRules []map[string]interface{}

	// If we return the API response directly, the output only gives
	// the contents of VirtualNetworkRuleProperties
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_server.listSQLServerVirtualNetworkRules", "api error", err)
		}
		for _, networkRule := range nextResult.Value {
			NetworkRules = append(NetworkRules, networkRuleMap(networkRule))
		}
	}

	return NetworkRules, nil
}

func networkRuleMap(rule *sql.VirtualNetworkRule) map[string]interface{} {
	objectMap := make(map[string]interface{})
	if rule.ID != nil {
		objectMap["id"] = rule.ID
	}
	if rule.Name != nil {
		objectMap["name"] = rule.Name
	}
	if rule.Type != nil {
		objectMap["type"] = rule.Type
	}
	if rule.Properties != nil {
		objectMap["properties"] = rule.Properties
	}
	return objectMap
}

// // If we return the API response directly, the output will not give
// // all the contents of PrivateEndpointConnection
func privateEndpointConnectionMap(conn *sql.PrivateEndpointConnection) PrivateConnectionInfo {
	var connection PrivateConnectionInfo
	if conn.ID != nil {
		connection.PrivateEndpointConnectionId = *conn.ID
	}
	if conn.Name != nil {
		connection.PrivateEndpointConnectionName = *conn.Name
	}
	if conn.Type != nil {
		connection.PrivateEndpointConnectionType = *conn.Type
	}
	if conn.Properties != nil {
		if conn.Properties.PrivateEndpoint != nil {
			if conn.Properties.PrivateEndpoint.ID != nil {
				connection.PrivateEndpointId = *conn.Properties.PrivateEndpoint.ID
			}
		}
		if conn.Properties.PrivateLinkServiceConnectionState != nil {
			if conn.Properties.PrivateLinkServiceConnectionState.ActionsRequired != nil {
				connection.PrivateLinkServiceConnectionStateActionsRequired = string(*conn.Properties.PrivateLinkServiceConnectionState.ActionsRequired)
			}
			if conn.Properties.PrivateLinkServiceConnectionState.Status != nil {
				connection.PrivateLinkServiceConnectionStateStatus = string(*conn.Properties.PrivateLinkServiceConnectionState.Status)
			}
			if conn.Properties.PrivateLinkServiceConnectionState.Description != nil {
				connection.PrivateLinkServiceConnectionStateDescription = *conn.Properties.PrivateLinkServiceConnectionState.Description
			}
		}
		if conn.Properties.ProvisioningState != nil {
			connection.ProvisioningState = string(*conn.Properties.ProvisioningState)
		}
	}

	return connection
}
