package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	sqlv "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureSQLServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_sql_server",
		Description: "Azure SQL Server",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getServer,
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
				Name:        "kind",
				Description: "The Kind of sql server",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "Location",
				Description: "The resource location",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrator_login",
				Description: "The Administrator username for the server",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AdministratorLogin"),
			},
			{
				Name:        "administrator_login_password",
				Description: "The administrator login password",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.AdministratorLoginPassword"),
			},
			{
				Name:        "version",
				Description: "The version of the server",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.Version"),
			},
			{
				Name:        "state",
				Description: "The state of the server",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.State"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the server",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerProperties.FullyQualifiedDomainName"),
			},
			{
				Name:        "audit_policy",
				Description: "Audit policies of server",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAuditPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_security_alert_policy",
				Description: "Server security alert policy",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServerSecurityAlertPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "server_azure_ad_administrator",
				Description: "Server Active Directory Administrator",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServerAzureADAdministrator,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "firewall_rules",
				Description: "he list of server firewall rules",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFirewallRules,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption_protector",
				Description: "The server encryption protector",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEncryptionProtector,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "Tags Attached to server",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
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

//// LIST FUNCTIONS ////

func listSQLServer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	serverClient := sql.NewServersClient(subscriptionID)
	serverClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := serverClient.List(context.Background())
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

//// HYDRATE FUNCTIONS ////

func getServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServer")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	serverClient := sql.NewServersClient(subscriptionID)
	serverClient.Authorizer = session.Authorizer

	op, err := serverClient.Get(ctx, resourceGroup, name)
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

func getFirewallRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFirewallRules")

	server := h.Item.(sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	FirewallRulesClient := sql.NewFirewallRulesClient(subscriptionID)
	FirewallRulesClient.Authorizer = session.Authorizer

	op, err := FirewallRulesClient.ListByServer(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getEncryptionProtector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEncryptionProtector")

	server := h.Item.(sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	EncryptionProtectorsClient := sql.NewEncryptionProtectorsClient(subscriptionID)
	EncryptionProtectorsClient.Authorizer = session.Authorizer

	op, err := EncryptionProtectorsClient.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getServerSecurityAlertPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServerSecurityAlertPolicy")

	server := h.Item.(sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	ServerSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClient(subscriptionID)
	ServerSecurityAlertPoliciesClient.Authorizer = session.Authorizer

	op, err := ServerSecurityAlertPoliciesClient.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getServerAzureADAdministrator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServerAzureADAdministrator")

	server := h.Item.(sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	ServerAzureADAdministratorClient := sql.NewServerAzureADAdministratorsClient(subscriptionID)
	ServerAzureADAdministratorClient.Authorizer = session.Authorizer

	op, err := ServerAzureADAdministratorClient.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		return nil, err
	}

	return op, nil
}

func getAuditPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditPolicy")

	server := h.Item.(sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	AuditPolicyClient := sql.NewServerBlobAuditingPoliciesClient(subscriptionID)
	AuditPolicyClient.Authorizer = session.Authorizer

	op, err := AuditPolicyClient.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getServerVulnerabilityAssessment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getServerVulnerabilityAssessment")

	server := h.Item.(sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	ServerVulnerabilityAssessmentClient := sqlv.NewServerVulnerabilityAssessmentsClient(subscriptionID)
	ServerVulnerabilityAssessmentClient.Authorizer = session.Authorizer

	op, err := ServerVulnerabilityAssessmentClient.Get(ctx, resourceGroupName, *server.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}