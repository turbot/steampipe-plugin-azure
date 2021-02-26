package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	sqlv "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
)

//// TABLE DEFINITION

func tableAzureSQLServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_sql_server",
		Description: "Azure SQL Server",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getSQLServer,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listSQLServer,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the SQS server",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a SQL server uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the SQS server",
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
				Description: "The Kind of sql server",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The resource location",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrator_login",
				Description: "Specifies the username of the Administrator for this server.",
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
				Description: "Specifies the active directory administrator",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLServerAzureADAdministrator,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_vulnerability_assessment",
				Description: "Server Active Directory Administrator",
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
				Name:        "tags_src",
				Description: "Specifies the set of tags attached to the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// steampipe standard columns
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
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// azure standard columns
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

func listSQLServer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewServersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := client.List(context.Background())
		if err != nil {
			return nil, err
		}
		for _, server := range result.Values() {
			d.StreamListItem(ctx, server)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getSQLServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServer")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewServersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	// // In some cases resource does not give any notFound error
	// // instead of notFound error, it returns empty data
	// if op.ID != nil {
	// 	return op, nil
	// }

	return op, nil
}

func getSQLServerAuditPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerAuditPolicy")

	server := h.Item.(sql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewServerBlobAuditingPoliciesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	auditPolicyData := map[string]interface{}{
		"id":         op.ID,
		"name":       op.Name,
		"type":       op.Type,
		"properties": op.ServerBlobAuditingPolicyProperties,
	}
	return auditPolicyData, nil
}

func getSQLServerSecurityAlertPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerSecurityAlertPolicy")

	server := h.Item.(sql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewServerSecurityAlertPoliciesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	securityAlertPolicyData := map[string]interface{}{
		"id":         op.ID,
		"name":       op.Name,
		"type":       op.Type,
		"properties": op.SecurityAlertPolicyProperties,
	}
	return securityAlertPolicyData, nil
}

func getSQLServerAzureADAdministrator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerAzureADAdministrator")

	server := h.Item.(sql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	serverAzureADAdministratorClient := sql.NewServerAzureADAdministratorsClient(subscriptionID)
	serverAzureADAdministratorClient.Authorizer = session.Authorizer

	op, err := serverAzureADAdministratorClient.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		return nil, err
	}

	administratorData := map[string]interface{}{
		"id":         op.ID,
		"name":       op.Name,
		"type":       op.Type,
		"properties": op.ServerAdministratorProperties,
	}
	return administratorData, nil
}

func getSQLServerEncryptionProtector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerEncryptionProtector")

	server := h.Item.(sql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewEncryptionProtectorsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	encryptionData := map[string]interface{}{
		"id":            op.ID,
		"name":          op.Name,
		"kind":          op.Kind,
		"subRegion":     op.EncryptionProtectorProperties.Subregion,
		"serverKeyName": op.EncryptionProtectorProperties.ServerKeyName,
		"serverKeyType": string(op.ServerKeyType),
		"uri":           op.URI,
		"thumbPrint":    op.Thumbprint,
	}

	return encryptionData, nil
}

func getSQLServerVulnerabilityAssessment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLServerVulnerabilityAssessment")

	server := h.Item.(sql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sqlv.NewServerVulnerabilityAssessmentsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	vulnerabilityAssessmentData := map[string]interface{}{
		"id":         op.ID,
		"name":       op.Name,
		"type":       op.Type,
		"properties": op.ServerVulnerabilityAssessmentProperties,
	}
	return vulnerabilityAssessmentData, nil
}

func listSQLServerFirewallRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLServerFirewallRules")

	server := h.Item.(sql.Server)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	client := sql.NewFirewallRulesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	if op.Value != nil {
		return op.Value, nil
	}

	return nil, nil
}
