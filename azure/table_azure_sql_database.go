package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	sql "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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
		Columns: azureColumns([]*plugin.Column{
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
				Name:        "requested_service_objective_name",
				Description: "The name of the configured service level objective of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.RequestedServiceObjectiveName"),
			},
			{
				Name:        "retention_policy_id",
				Description: "Retention policy ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSqlDatabaseLongTermRetentionPolicies,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "retention_policy_name",
				Description: "Retention policy Name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSqlDatabaseLongTermRetentionPolicies,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "retention_policy_type",
				Description: "Long term Retention policy Type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSqlDatabaseLongTermRetentionPolicies,
				Transform:   transform.FromField("Type"),
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
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.CreateMode"),
			},
			{
				Name:        "read_scale",
				Description: "ReadScale indicates whether read-only connections are allowed to this database or not if the database is a geo-secondary.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseProperties.ReadScale"),
			},
			{
				Name:        "recommended_index",
				Description: "The recommended indices for this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DatabaseProperties.RecommendedIndex"),
			},
			{
				Name:        "retention_policy_property",
				Description: "Long term Retention policy Property.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSqlDatabaseLongTermRetentionPolicies,
				Transform:   transform.FromField("BaseLongTermRetentionPolicyProperties"),
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
			{
				Name:        "vulnerability_assessments",
				Description: "The vulnerability assessments for this database.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSqlDatabaseVulnerabilityAssessments,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "vulnerability_assessment_scan_records",
				Description: "The vulnerability assessment scan records for this database.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSqlDatabaseVulnerabilityAssessmentScans,
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
				Hydrate:     getSqlDatabase,
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

func listSqlDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabases", "connection error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewDatabasesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabases", "client error", err)
	}

	server := h.Item.(*sql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	pager := client.NewListByServerPager(resourceGroupName, *server.Name, nil)

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabases", "api error", err)
		}
		for _, database := range nextResult.Value {
			d.StreamListItem(ctx, database)
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

func getSqlDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var serverName, databaseName, resourceGroupName string
	if h.Item != nil {
		database := h.Item.(*sql.Database)
		serverName = strings.Split(*database.ID, "/")[8]
		databaseName = *database.Name
		resourceGroupName = strings.Split(string(*database.ID), "/")[4]
	} else {
		serverName = d.KeyColumnQuals["server_name"].GetStringValue()
		databaseName = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroupName = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabase", "credential error", err)
	}
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewDatabasesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabase", "client error", err)
	}

	op, err := client.Get(ctx, resourceGroupName, serverName, databaseName, nil)
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
	database := h.Item.(*sql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	databaseName := *database.Name
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabaseTransparentDataEncryption", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewTransparentDataEncryptionsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabaseTransparentDataEncryption", "client error", err)
	}

	op := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, nil)

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabaseTransparentDataEncryption", "api error", err)
		}

		if len(nextResult.LogicalDatabaseTransparentDataEncryptionListResult.Value) > 0 {
			return nextResult.LogicalDatabaseTransparentDataEncryptionListResult.Value[0], nil
		}
	}

	return nil, nil
}

func getSqlDatabaseLongTermRetentionPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(*sql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	databaseName := *database.Name
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabaseLongTermRetentionPolicies", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewLongTermRetentionPoliciesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabaseLongTermRetentionPolicies", "client error", err)
	}

	op := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, nil)
	if err != nil {
		return nil, err
	}
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_database.getSqlDatabaseLongTermRetentionPolicies", "api error", err)
		}
		// We can add only one retention policy per SQL Database.
		if len(nextResult.LongTermRetentionPolicyListResult.Value) > 0 {
			return nextResult.LongTermRetentionPolicyListResult.Value[0], nil
		}
	}

	return nil, nil
}

func listSqlDatabaseVulnerabilityAssessments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(*sql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	databaseName := *database.Name
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabaseVulnerabilityAssessments", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewDatabaseVulnerabilityAssessmentsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabaseVulnerabilityAssessments", "client error", err)
	}

	op := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	var vulnerabilityAssessments []map[string]interface{}
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabaseVulnerabilityAssessments", "api error", err)
		}
		objectMap := make(map[string]interface{})
		for _, i := range nextResult.Value {
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.Name != nil {
				objectMap["name"] = i.Name
			}
			if i.Type != nil {
				objectMap["type"] = i.Type
			}
			if i.Properties.RecurringScans != nil {
				objectMap["recurringScans"] = i.Properties.RecurringScans
			}
			if i.Properties.StorageAccountAccessKey != nil {
				objectMap["storageAccountAccessKey"] = *i.Properties.StorageAccountAccessKey
			}
			if i.Properties.StorageContainerPath != nil {
				objectMap["storageContainerPath"] = *i.Properties.StorageContainerPath
			}
			if i.Properties.StorageContainerSasKey != nil {
				objectMap["storageContainerSasKey"] = *i.Properties.StorageContainerSasKey
			}
			vulnerabilityAssessments = append(vulnerabilityAssessments, objectMap)
		}
	}

	return vulnerabilityAssessments, nil
}

func listSqlDatabaseVulnerabilityAssessmentScans(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(*sql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	databaseName := *database.Name
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabaseVulnerabilityAssessmentScans", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewDatabaseVulnerabilityAssessmentScansClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabaseVulnerabilityAssessmentScans", "client error", err)
	}

	op := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, sql.VulnerabilityAssessmentNameDefault, nil)
	var vulnerabilityAssessmentScanRecords []map[string]interface{}
	if err != nil {
		// API throws "VulnerabilityAssessmentInvalidPolicy" error if Vulnerability Assessment settings don't exist or invalid storage specified in settings.
		// https://learn.microsoft.com/en-us/rest/api/sql/2022-05-01-preview/database-vulnerability-assessment-scans/list-by-database?tabs=HTTP
		if strings.Contains(err.Error(), "VulnerabilityAssessmentInvalidPolicy") {
			return vulnerabilityAssessmentScanRecords, nil
		}
		return nil, err
	}

	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_sql_database.listSqlDatabaseVulnerabilityAssessmentScans", "api error", err)
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
			if i.Properties.ScanID != nil {
				objectMap["scanID"] = *i.Properties.ScanID
			}
			if len(*i.Properties.TriggerType) > 0 {
				objectMap["triggerType"] = i.Properties.TriggerType
			}
			if len(*i.Properties.State) > 0 {
				objectMap["state"] = i.Properties.State
			}
			if i.Properties.StartTime != nil {
				objectMap["startTime"] = i.Properties.StartTime
			}
			if i.Properties.EndTime != nil {
				objectMap["endTime"] = i.Properties.EndTime
			}
			if i.Properties.Errors != nil {
				objectMap["errors"] = i.Properties.Errors
			}
			if i.Properties.StorageContainerPath != nil {
				objectMap["storageContainerPath"] = i.Properties.StorageContainerPath
			}
			if i.Properties.NumberOfFailedSecurityChecks != nil {
				objectMap["numberOfFailedSecurityChecks"] = *i.Properties.NumberOfFailedSecurityChecks
			}
			vulnerabilityAssessmentScanRecords = append(vulnerabilityAssessmentScanRecords, objectMap)
		}
	}

	return vulnerabilityAssessmentScanRecords, nil
}

//// TRANSFORM FUNCTION

func idToServerName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(sql.Database)
	serverName := strings.Split(string(*data.ID), "/")[8]
	return serverName, nil
}
