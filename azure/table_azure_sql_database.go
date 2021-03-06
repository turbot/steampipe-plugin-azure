package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSqlDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_sql_database",
		Description: "Azure SQL Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "server_name", "resource_group"}),
			Hydrate:    getSqlDatabase,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listSQLServer,
			Hydrate:       listSqlDatabases,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the database.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a database uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "server_name",
				Description: "The name of the parent server of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(idToServerName),
			},
			{
				Name:        "status",
				Description: "The status of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.Status"),
			},
			{
				Name:        "type",
				Description: "Type of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "collation",
				Description: "The collation of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.Collation"),
			},
			{
				Name:        "containment_state",
				Description: "The containment state of the database.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DatabaseProperties.ContainmentState"),
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DatabaseProperties.CreationDate").Transform(convertDateToTime),
			},
			{
				Name:        "current_service_objective_id",
				Description: "The current service level objective ID of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.CurrentServiceObjectiveID"),
			},
			{
				Name:        "database_id",
				Description: "The ID of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.DatabaseID"),
			},
			{
				Name:        "default_secondary_location",
				Description: "The default secondary region for this database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.DefaultSecondaryLocation"),
			},
			{
				Name:        "earliest_restore_date",
				Description: "This records the earliest start date and time that restore is available for this database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DatabaseProperties.EarliestRestoreDate").Transform(convertDateToTime),
			},
			{
				Name:        "edition",
				Description: "The edition of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.Edition"),
			},
			{
				Name:        "elastic_pool_name",
				Description: "The name of the elastic pool the database is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.ElasticPoolName"),
			},
			{
				Name:        "failover_group_id",
				Description: "The resource identifier of the failover group containing this database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.FailoverGroupID"),
			},
			{
				Name:        "kind",
				Description: "Kind of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "Location of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_size_bytes",
				Description: "The max size of the database expressed in bytes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.MaxSizeBytes"),
			},
			{
				Name:        "recovery_services_recovery_point_resource_id",
				Description: "Specifies the resource ID of the recovery point to restore from if createMode is RestoreLongTermRetentionBackup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.RecoveryServicesRecoveryPointResourceID"),
			},
			{
				Name:        "requested_service_objective_id",
				Description: "The configured service level objective ID of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.RequestedServiceObjectiveID"),
			},
			{
				Name:        "restore_point_in_time",
				Description: "Specifies the point in time of the source database that will be restored to create the new database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DatabaseProperties.RestorePointInTime").Transform(convertDateToTime),
			},
			{
				Name:        "source_database_deletion_date",
				Description: "Specifies the time that the database was deleted when createMode is Restore and sourceDatabaseId is the deleted database's original resource id.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DatabaseProperties.SourceDatabaseDeletionDate").Transform(convertDateToTime),
			},
			{
				Name:        "source_database_id",
				Description: "Specifies the resource ID of the source database if createMode is Copy, NonReadableSecondary, OnlineSecondary, PointInTimeRestore, Recovery, or Restore.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.SourceDatabaseID"),
			},
			{
				Name:        "zone_redundant",
				Description: "Indicates if the database is zone redundant or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DatabaseProperties.ZoneRedundant"),
			},
			{
				Name:        "create_mode",
				Description: "Specifies the mode of database creation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.CreateMode"),
			},
			{
				Name:        "read_scale",
				Description: "ReadScale indicates whether read-only connections are allowed to this database or not if the database is a geo-secondary.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.ReadScale"),
			},
			{
				Name:        "recommended_index",
				Description: "The recommended indices for this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.RecommendedIndex"),
			},
			{
				Name:        "requested_service_objective_name",
				Description: "The name of the configured service level objective of the database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.RequestedServiceObjectiveName"),
			},
			{
				Name:        "sample_name",
				Description: "Indicates the name of the sample schema to apply when creating this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.SampleName"),
			},
			{
				Name:        "service_level_objective",
				Description: "The current service level objective of the database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.ServiceLevelObjective"),
			},
			{
				Name:        "service_tier_advisors",
				Description: "The list of service tier advisors for this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.ServiceTierAdvisors"),
			},
			{
				Name:        "transparent_data_encryption",
				Description: "The transparent data encryption info for this database.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSqlDatabaseTransparentDataEncryption,
				Transform:   transform.FromField("TransparentDataEncryptionProperties"),
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
				Hydrate:     getSqlDatabase,
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

//// LIST FUNCTION

func listSqlDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewDatabasesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	server := h.Item.(sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	result, err := client.ListByServer(ctx, resourceGroupName, *server.Name, "", "")
	if err != nil {
		return nil, err
	}

	for _, database := range *result.Value {
		d.StreamLeafListItem(ctx, database)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSqlDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSqlDatabase")

	var serverName, databaseName, resourceGroupName string
	if h.Item != nil {
		database := h.Item.(sql.Database)
		serverName = strings.Split(*database.ID, "/")[8]
		databaseName = *database.Name
		resourceGroupName = strings.Split(string(*database.ID), "/")[4]
	} else {
		serverName = d.KeyColumnQuals["server_name"].GetStringValue()
		databaseName = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroupName = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewDatabasesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, serverName, databaseName, "")
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

func getSqlDatabaseTransparentDataEncryption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSqlDatabaseTransparentDataEncryption")

	var serverName, databaseName, resourceGroupName string
	if h.Item != nil {
		database := h.Item.(sql.Database)
		serverName = strings.Split(*database.ID, "/")[8]
		databaseName = *database.Name
		resourceGroupName = strings.Split(string(*database.ID), "/")[4]
	} else {
		serverName = d.KeyColumnQuals["server_name"].GetStringValue()
		databaseName = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroupName = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewTransparentDataEncryptionsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroupName, serverName, databaseName)
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

//// TRANSFORM FUNCTION

func idToServerName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(sql.Database)
	serverName := strings.Split(string(*data.ID), "/")[8]
	return serverName, nil
}
