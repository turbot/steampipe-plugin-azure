package azure

import (
	"context"
	"reflect"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	armmypostgresflexibleservers "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v2"
)

//// TABLE DEFINITION

func tableAzurePostgreSqlFlexibleServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_postgresql_flexible_server",
		Description: "Azure PostgreSQL Flexible Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getPostgreSqlFlexibleServer,
			Tags: map[string]string{
				"service": "Microsoft.DBforPostgreSQL",
				"action":  "flexibleServers/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listPostgreSqlFlexibleServers,
			Tags: map[string]string{
				"service": "Microsoft.DBforPostgreSQL",
				"action":  "flexibleServers/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "A state of a server that is visible to user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.State"),
			},
			{
				Name:        "availability_zone",
				Description: "Availability zone information of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AvailabilityZone"),
			},
			{
				Name:        "administrator_login",
				Description: "The administrator's login name of a server. Can only be specified when the server is being created (and is required for creation).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AdministratorLogin"),
			},
			{
				Name:        "create_mode",
				Description: "The mode to create a new PostgreSQL server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CreateMode"),
			},
			{
				Name:        "point_in_time_utc",
				Description: "Restore point creation time (ISO8601 format), specifying the time to restore from. It's required when 'createMode' is 'PointInTimeRestore' or 'GeoRestore'.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.PointInTimeUTC").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "replica_capacity",
				Description: "Replicas allowed for a server.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.ReplicaCapacity"),
			},
			{
				Name:        "replication_role",
				Description: "Replication role of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ReplicationRole"),
			},
			{
				Name:        "source_server_resource_id",
				Description: "The source server resource ID to restore from. It's required when 'createMode' is 'PointInTimeRestore' or 'GeoRestore' or 'Replica'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SourceServerResourceID"),
			},
			{
				Name:        "version",
				Description: "PostgreSQL Server version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Version"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of a server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.FullyQualifiedDomainName"),
			},
			{
				Name:        "minor_version",
				Description: "The minor version of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinorVersion"),
			},
			{
				Name:        "public_network_access",
				Description: "Public network access is enabled or not.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Network.PublicNetworkAccess"),
			},
			{
				Name:        "location",
				Description: "The geo-location where the resource lives.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auth_config",
				Description: "AuthConfig properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AuthConfig"),
			},
			{
				Name:        "backup",
				Description: "Backup properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Backup"),
			},
			{
				Name:        "data_encryption",
				Description: "Data encryption properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.DataEncryption"),
			},
			{
				Name:        "high_availability",
				Description: "High availability properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.HighAvailability"),
			},
			{
				Name:        "maintenance_window",
				Description: "Maintenance window properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.MaintenanceWindow"),
			},
			{
				Name:        "network",
				Description: "Network properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Network"),
			},
			{
				Name:        "storage",
				Description: "Storage properties of a server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Storage"),
			},
			{
				Name:        "sku",
				Description: "The SKU (pricing tier) of the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SKU"),
			},
			{
				Name:        "system_data",
				Description: "The system metadata relating to this server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "server_properties",
				Description: "Properties of the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties").Transform(extractPostgresFlexibleServerProperties),
			},
			{
				Name:        "flexible_server_configurations",
				Description: "The server configurations(parameters) details of the server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listPostgreSQLFlexibleServersConfigurations,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "firewall_rules",
				Description: "The list of firewall rules in a server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listPostgreSQLFlexibleServerFirewallRules,
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

func listPostgreSqlFlexibleServers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := armmypostgresflexibleservers.NewServersClient(subscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSqlFlexibleServers", "session_error", err)
		return nil, err
	}
	resourceGroupName := h.Item.(resources.Group).Name

	input := &armmypostgresflexibleservers.ServersClientListByResourceGroupOptions{}

	pager := client.NewListByResourceGroupPager(*resourceGroupName, input)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSqlFlexibleServers", "api_error", err)
			return nil, err
		}
		for _, server := range page.Value {
			d.StreamListItem(ctx, *server)

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

func getPostgreSqlFlexibleServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQualString("name")
	resourceGroup := d.EqualsQualString("resource_group")

	// check if name or resourceGroup is empty
	if resourceGroup == "" || name == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := armmypostgresflexibleservers.NewServersClient(subscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.getPostgreSqlFlexibleServer", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.getPostgreSqlFlexibleServer", "api_error", err)
		return nil, err
	}

	return op.Server, nil
}

func listPostgreSQLFlexibleServersConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	server := h.Item.(armmypostgresflexibleservers.Server)
	resourceGroup := strings.Split(string(*server.ID), "/")[4]
	serverName := *server.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := armmypostgresflexibleservers.NewConfigurationsClient(subscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSQLFlexibleServersConfigurations", "client_error", err)
		return nil, err
	}

	var postgreSQLFlexibleServersConfigurations []map[string]interface{}

	input := &armmypostgresflexibleservers.ConfigurationsClientListByServerOptions{}
	pager := client.NewListByServerPager(resourceGroup, serverName, input)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSQLFlexibleServersConfigurations", "api_error", err)
			return nil, err
		}

		for _, configuration := range page.Value {
			postgreSQLFlexibleServersConfigurations = append(postgreSQLFlexibleServersConfigurations, extractpostgreSQLFlexibleServersconfiguration(*configuration))
		}
	}

	return postgreSQLFlexibleServersConfigurations, nil
}

func listPostgreSQLFlexibleServerFirewallRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	server := h.Item.(armmypostgresflexibleservers.Server)
	resourceGroup := strings.Split(string(*server.ID), "/")[4]
	serverName := *server.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := armmypostgresflexibleservers.NewFirewallRulesClient(subscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSQLFlexibleServerFirewallRules", "client_error", err)
		return nil, err
	}

	var postgreSQLFlexibleServerFirewallRules []*armmypostgresflexibleservers.FirewallRule

	input := &armmypostgresflexibleservers.FirewallRulesClientListByServerOptions{}
	pager := client.NewListByServerPager(resourceGroup, serverName, input)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSQLFlexibleServerFirewallRules", "api_error", err)
			return nil, err
		}
		postgreSQLFlexibleServerFirewallRules = append(postgreSQLFlexibleServerFirewallRules, page.Value...)
	}

	return postgreSQLFlexibleServerFirewallRules, nil
}

//// TRANSFORM FUNCTION

// If we return the API response directly, the output will not provide the properties of Configurations
func extractpostgreSQLFlexibleServersconfiguration(i armmypostgresflexibleservers.Configuration) map[string]interface{} {
	postgreSQLFlexibleServersconfiguration := make(map[string]interface{})

	if i.ID != nil {
		postgreSQLFlexibleServersconfiguration["ID"] = *i.ID
	}
	if i.Name != nil {
		postgreSQLFlexibleServersconfiguration["Name"] = *i.Name
	}
	if i.Type != nil {
		postgreSQLFlexibleServersconfiguration["Type"] = *i.Type
	}
	if i.Properties != nil {
		postgreSQLFlexibleServersconfiguration["ConfigurationProperties"] = *i.Properties
	}

	return postgreSQLFlexibleServersconfiguration
}

func extractPostgresFlexibleServerProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	conf := d.HydrateItem.(armmypostgresflexibleservers.Server)
	if conf.Properties != nil {
		return structToMap(reflect.ValueOf(*conf.Properties)), nil
	}
	return nil, nil
}
