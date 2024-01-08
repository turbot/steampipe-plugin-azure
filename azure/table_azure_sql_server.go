package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	sqlv3 "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
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
				Transform:   transform.FromField("ServerProperties.State"),
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
				Transform:   transform.FromField("ServerProperties.AdministratorLogin"),
			},
			{
				Name:        "administrator_login_password",
				Description: "The administrator login password.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AdministratorLoginPassword"),
			},
			{
				Name:        "minimal_tls_version",
				Description: "Minimal TLS version. Allowed values: '1.0', '1.1', '1.2'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.MinimalTLSVersion"),
			},
			{
				Name:        "public_network_access",
				Description: "Whether or not public endpoint access is allowed for this server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.PublicNetworkAccess"),
			},
			{
				Name:        "version",
				Description: "The version of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Version"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.FullyQualifiedDomainName"),
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
			{
				Name:        "audit_policy",
				Description: "The SQL server blob auditing policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLServerBlobAuditingPolicies,
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
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := sqlv3.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, server := range result.Values() {
		d.StreamListItem(ctx, server)
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
		for _, server := range result.Values() {
			d.StreamListItem(ctx, server)
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

func getSQLServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServer")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sqlv3.NewServersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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

func getSQLServerAuditPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerAuditPolicy")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewServerBlobAuditingPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of ServerBlobAuditingPolicyProperties
	var auditPolicies []map[string]interface{}
	for _, i := range op.Values() {
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
		if i.ServerBlobAuditingPolicyProperties != nil {
			objectMap["properties"] = i.ServerBlobAuditingPolicyProperties
		}
		auditPolicies = append(auditPolicies, objectMap)
	}
	return auditPolicies, nil
}

func listSQLServerPrivateEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLServerPrivateEndpointConnections")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sqlv3.NewPrivateEndpointConnectionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		plugin.Logger(ctx).Error("listSQLServerPrivateEndpointConnections", "ListByServer", err)
		return nil, err
	}

	var privateEndpointConnections []PrivateConnectionInfo
	var connection PrivateConnectionInfo

	for _, conn := range op.Values() {
		connection = privateEndpointConnectionMap(conn)
		privateEndpointConnections = append(privateEndpointConnections, connection)
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listSQLServerPrivateEndpointConnections", "ListByServer_pagination", err)
			return nil, err
		}
		for _, conn := range op.Values() {
			connection = privateEndpointConnectionMap(conn)
			privateEndpointConnections = append(privateEndpointConnections, connection)
		}
	}
	return privateEndpointConnections, nil
}

func getSQLServerSecurityAlertPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerSecurityAlertPolicy")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of SecurityAlertPolicyProperties
	var securityAlertPolicies []map[string]interface{}
	for _, i := range op.Values() {
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
		if i.SecurityAlertPolicyProperties != nil {
			objectMap["properties"] = i.SecurityAlertPolicyProperties
		}
		securityAlertPolicies = append(securityAlertPolicies, objectMap)
	}
	return securityAlertPolicies, nil
}

func getSQLServerBlobAuditingPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewServerBlobAuditingPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	var blobPolicies []map[string]interface{}
	for _, i := range op.Values() {
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
		if i.ServerBlobAuditingPolicyProperties != nil {
			obMap := make(map[string]interface{})
			if i.ServerBlobAuditingPolicyProperties.RetentionDays != nil {
				obMap["retentionDays"] = i.ServerBlobAuditingPolicyProperties.RetentionDays
			}
			if i.ServerBlobAuditingPolicyProperties.AuditActionsAndGroups != nil {
				obMap["AuditActionsAndGroups"] = i.ServerBlobAuditingPolicyProperties.AuditActionsAndGroups
			}
			if i.ServerBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled != nil {
				obMap["isAzureMonitorTargetEnabled"] = i.ServerBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled
			}
			if i.ServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse != nil {
				obMap["isStorageSecondaryKeyInUse"] = i.ServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse
			}
			if i.ServerBlobAuditingPolicyProperties.QueueDelayMs != nil {
				obMap["queueDelayMs"] = i.ServerBlobAuditingPolicyProperties.QueueDelayMs
			}
			if i.ServerBlobAuditingPolicyProperties.State != "" {
				obMap["state"] = i.ServerBlobAuditingPolicyProperties.State
			}
			if i.ServerBlobAuditingPolicyProperties.StorageEndpoint != nil {
				obMap["storageEndpoint"] = i.ServerBlobAuditingPolicyProperties.StorageEndpoint
			}
			if i.ServerBlobAuditingPolicyProperties.StorageAccountAccessKey != nil {
				obMap["storageAccountAccessKey"] = i.ServerBlobAuditingPolicyProperties.StorageAccountAccessKey
			}
			if i.ServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID != nil {
				obMap["storageAccountSubscriptionID"] = i.ServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID
			}
			objectMap["serverBlobAuditingPolicyProperties"] = obMap
		}

		blobPolicies = append(blobPolicies, objectMap)
	}

	if op.NotDone() {
		for _, i := range op.Values() {
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
			if i.ServerBlobAuditingPolicyProperties != nil {
				obMap := make(map[string]interface{})
				if i.ServerBlobAuditingPolicyProperties.RetentionDays != nil {
					obMap["retentionDays"] = i.ServerBlobAuditingPolicyProperties.RetentionDays
				}
				if i.ServerBlobAuditingPolicyProperties.AuditActionsAndGroups != nil {
					obMap["AuditActionsAndGroups"] = i.ServerBlobAuditingPolicyProperties.AuditActionsAndGroups
				}
				if i.ServerBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled != nil {
					obMap["isAzureMonitorTargetEnabled"] = i.ServerBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled
				}
				if i.ServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse != nil {
					obMap["isStorageSecondaryKeyInUse"] = i.ServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse
				}
				if i.ServerBlobAuditingPolicyProperties.QueueDelayMs != nil {
					obMap["queueDelayMs"] = i.ServerBlobAuditingPolicyProperties.QueueDelayMs
				}
				if i.ServerBlobAuditingPolicyProperties.State != "" {
					obMap["state"] = i.ServerBlobAuditingPolicyProperties.State
				}
				if i.ServerBlobAuditingPolicyProperties.StorageEndpoint != nil {
					obMap["storageEndpoint"] = i.ServerBlobAuditingPolicyProperties.StorageEndpoint
				}
				if i.ServerBlobAuditingPolicyProperties.StorageAccountAccessKey != nil {
					obMap["storageAccountAccessKey"] = i.ServerBlobAuditingPolicyProperties.StorageAccountAccessKey
				}
				if i.ServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID != nil {
					obMap["storageAccountSubscriptionID"] = i.ServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID
				}
				objectMap["serverBlobAuditingPolicyProperties"] = obMap
			}

			blobPolicies = append(blobPolicies, objectMap)
		}
	}

	return blobPolicies, nil
}

func getSQLServerAzureADAdministrator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerAzureADAdministrator")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewServerAzureADAdministratorsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of ServerAdministratorProperties
	var serverAdministrators []map[string]interface{}
	for _, i := range *op.Value {
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
		if i.ServerAdministratorProperties != nil {
			objectMap["properties"] = i.ServerAdministratorProperties
		}
		serverAdministrators = append(serverAdministrators, objectMap)
	}
	return serverAdministrators, nil
}

func getSQLServerEncryptionProtector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerEncryptionProtector")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewEncryptionProtectorsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of EncryptionProtectorProperties
	var encryptionProtectors []map[string]interface{}
	for _, i := range op.Values() {
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
		if i.Subregion != nil {
			objectMap["subregion"] = i.Subregion
		}
		if i.ServerKeyName != nil {
			objectMap["serverKeyName"] = i.ServerKeyName
		}
		if i.ServerKeyType != "" {
			objectMap["serverKeyType"] = i.ServerKeyType
		}
		if i.URI != nil {
			objectMap["uri"] = i.URI
		}
		if i.Thumbprint != nil {
			objectMap["thumbprint"] = i.Thumbprint
		}
		encryptionProtectors = append(encryptionProtectors, objectMap)
	}

	return encryptionProtectors, nil
}

func getSQLServerVulnerabilityAssessment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerVulnerabilityAssessment")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sqlv3.NewServerVulnerabilityAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of ServerVulnerabilityAssessmentProperties
	var vulnerabilityAssessments []map[string]interface{}
	for _, i := range op.Values() {
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
		if i.ServerVulnerabilityAssessmentProperties != nil {
			objectMap["properties"] = i.ServerVulnerabilityAssessmentProperties
		}
		vulnerabilityAssessments = append(vulnerabilityAssessments, objectMap)
	}
	return vulnerabilityAssessments, nil
}

func listSQLServerFirewallRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLServerFirewallRules")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewFirewallRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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
			objectMap["id"] = i.ID
		}
		if i.Name != nil {
			objectMap["name"] = i.Name
		}
		if i.Type != nil {
			objectMap["type"] = i.Type
		}
		if i.FirewallRuleProperties != nil {
			objectMap["properties"] = i.FirewallRuleProperties
		}
		firewallRules = append(firewallRules, objectMap)
	}

	return firewallRules, nil
}

func listSQLServerVirtualNetworkRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLServerVirtualNetworkRules")
	server := h.Item.(sqlv3.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewVirtualNetworkRulesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// If we return the API response directly, the output only gives
	// the contents of VirtualNetworkRuleProperties
	var NetworkRules []map[string]interface{}
	result, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}
	for _, networkRule := range result.Values() {
		NetworkRules = append(NetworkRules, networkRuleMap(networkRule))
	}

	for result.NotDone() {
		err := result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, networkRule := range result.Values() {
			NetworkRules = append(NetworkRules, networkRuleMap(networkRule))
		}
	}

	return NetworkRules, nil
}

func networkRuleMap(rule sql.VirtualNetworkRule) map[string]interface{} {
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
	if rule.VirtualNetworkRuleProperties != nil {
		objectMap["properties"] = rule.VirtualNetworkRuleProperties
	}
	return objectMap
}

// If we return the API response directly, the output will not give
// all the contents of PrivateEndpointConnection
func privateEndpointConnectionMap(conn sqlv3.PrivateEndpointConnection) PrivateConnectionInfo {
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
	if conn.PrivateEndpointConnectionProperties != nil {
		if conn.PrivateEndpoint != nil {
			if conn.PrivateEndpoint.ID != nil {
				connection.PrivateEndpointId = *conn.PrivateEndpoint.ID
			}
		}
		if conn.PrivateLinkServiceConnectionState != nil {
			if conn.PrivateLinkServiceConnectionState.ActionsRequired != "" {
				connection.PrivateLinkServiceConnectionStateActionsRequired = string(conn.PrivateLinkServiceConnectionState.ActionsRequired)
			}
			if conn.PrivateLinkServiceConnectionState.Status != "" {
				connection.PrivateLinkServiceConnectionStateStatus = string(conn.PrivateLinkServiceConnectionState.Status)
			}
			if conn.PrivateLinkServiceConnectionState.Description != nil {
				connection.PrivateLinkServiceConnectionStateDescription = *conn.PrivateLinkServiceConnectionState.Description
			}
		}
		if conn.ProvisioningState != "" {
			connection.ProvisioningState = string(conn.ProvisioningState)
		}
	}

	return connection
}
