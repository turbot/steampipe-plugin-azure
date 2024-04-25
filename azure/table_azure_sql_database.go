package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:    listSqlDatabaseVulnerabilityAssessmentScans,
				Depends: []plugin.HydrateFunc{listSqlDatabaseVulnerabilityAssessments},
			},
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
				Transform:   transform.FromField("Properties.Status"),
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
				Transform:   transform.FromField("Properties.Collation"),
			},
			{
				Name:        "containment_state",
				Description: "The containment state of the database.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.ContainmentState"),
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.CreationDate"),
			},
			{
				Name:        "current_service_objective_id",
				Description: "The current service level objective ID of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CurrentServiceObjectiveID"),
			},
			{
				Name:        "database_id",
				Description: "The ID of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DatabaseID"),
			},
			{
				Name:        "default_secondary_location",
				Description: "The default secondary region for this database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DefaultSecondaryLocation"),
			},
			{
				Name:        "earliest_restore_date",
				Description: "This records the earliest start date and time that restore is available for this database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.EarliestRestoreDate"),
			},
			{
				Name:        "edition",
				Description: "The edition of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Edition"),
			},
			{
				Name:        "elastic_pool_name",
				Description: "The name of the elastic pool the database is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ElasticPoolName"),
			},
			{
				Name:        "failover_group_id",
				Description: "The resource identifier of the failover group containing this database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.FailoverGroupID"),
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
				Transform:   transform.FromField("Properties.MaxSizeBytes"),
			},
			{
				Name:        "recovery_services_recovery_point_resource_id",
				Description: "Specifies the resource ID of the recovery point to restore from if createMode is RestoreLongTermRetentionBackup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RecoveryServicesRecoveryPointResourceID"),
			},
			{
				Name:        "requested_service_objective_id",
				Description: "The configured service level objective ID of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RequestedServiceObjectiveID"),
			},
			{
				Name:        "restore_point_in_time",
				Description: "Specifies the point in time of the source database that will be restored to create the new database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.RestorePointInTime"),
			},
			{
				Name:        "requested_service_objective_name",
				Description: "The name of the configured service level objective of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RequestedServiceObjectiveName"),
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
				Transform:   transform.FromField("Properties.SourceDatabaseDeletionDate"),
			},
			{
				Name:        "source_database_id",
				Description: "Specifies the resource ID of the source database if createMode is Copy, NonReadableSecondary, OnlineSecondary, PointInTimeRestore, Recovery, or Restore.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SourceDatabaseID"),
			},
			{
				Name:        "zone_redundant",
				Description: "Indicates if the database is zone redundant or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.ZoneRedundant"),
			},
			{
				Name:        "create_mode",
				Description: "Specifies the mode of database creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CreateMode"),
			},
			{
				Name:        "read_scale",
				Description: "ReadScale indicates whether read-only connections are allowed to this database or not if the database is a geo-secondary.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ReadScale"),
			},
			{
				Name:        "recommended_index",
				Description: "The recommended indices for this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.RecommendedIndex"),
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
				Transform:   transform.FromField("Properties.ServiceLevelObjective"),
			},
			{
				Name:        "service_tier_advisors",
				Description: "The list of service tier advisors for this database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ServiceTierAdvisors"),
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
			{
				Name:        "audit_policy",
				Description: "The database blob auditing policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSqlDatabaseBlobAuditingPolicies,
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
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewDatabasesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	server := h.Item.(armsql.Server)
	resourceGroupName := strings.Split(string(*server.ID), "/")[4]

	pager := client.NewListByServerPager(resourceGroupName, *server.Name, nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, database := range result.Value {
			d.StreamListItem(ctx, *database)
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

func getSqlDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSqlDatabase")

	var serverName, databaseName, resourceGroupName string
	if h.Item != nil {
		database := h.Item.(armsql.Database)
		serverName = strings.Split(*database.ID, "/")[8]
		databaseName = *database.Name
		resourceGroupName = strings.Split(string(*database.ID), "/")[4]
	} else {
		serverName = d.EqualsQuals["server_name"].GetStringValue()
		databaseName = d.EqualsQuals["name"].GetStringValue()
		resourceGroupName = d.EqualsQuals["resource_group"].GetStringValue()
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewDatabasesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroupName, serverName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op.Database, nil
	}

	return nil, nil
}

func getSqlDatabaseTransparentDataEncryption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(armsql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]
	databaseName := *database.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewTransparentDataEncryptionsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	var tdes []*armsql.LogicalDatabaseTransparentDataEncryption
	pager := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		tdes = append(tdes, result.Value...)
	}

	return tdes, nil
}

func getSqlDatabaseLongTermRetentionPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(armsql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]
	databaseName := *database.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewLongTermRetentionPoliciesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	pager := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if len(result.Value) > 0 {
			// Return the first retention policy found
			return result.Value[0], nil
		}
	}

	return nil, nil
}

func getSqlDatabaseBlobAuditingPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(armsql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]
	databaseName := *database.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewDatabaseBlobAuditingPoliciesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	var blobPolicies []*armsql.DatabaseBlobAuditingPolicy
	pager := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		blobPolicies = append(blobPolicies, result.Value...)
	}

	return blobPolicies, nil
}

func listSqlDatabaseVulnerabilityAssessments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(armsql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]
	databaseName := *database.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewDatabaseVulnerabilityAssessmentsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	var vulnerabilityAssessments []*armsql.DatabaseVulnerabilityAssessment
	pager := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		vulnerabilityAssessments = append(vulnerabilityAssessments, result.Value...)
	}

	return vulnerabilityAssessments, nil
}

func listSqlDatabaseVulnerabilityAssessmentScans(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(armsql.Database)
	serverName := strings.Split(*database.ID, "/")[8]
	resourceGroupName := strings.Split(string(*database.ID), "/")[4]
	databaseName := *database.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewDatabaseVulnerabilityAssessmentScansClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}
	vulnerabilityAssessments := h.HydrateResults["listSqlDatabaseVulnerabilityAssessments"].([]*armsql.DatabaseVulnerabilityAssessment)
	var vulnerabilityAssessmentScanRecords []*armsql.VulnerabilityAssessmentScanRecord

	for _, vulnerabilityAssessment := range vulnerabilityAssessments {
		pager := client.NewListByDatabasePager(resourceGroupName, serverName, databaseName, armsql.VulnerabilityAssessmentName(*vulnerabilityAssessment.Name), nil)
		for pager.More() {
			result, err := pager.NextPage(ctx)
			if err != nil {

				// check if Vulnerability Assessment is invalid
				if strings.Contains(err.Error(), "VulnerabilityAssessmentInvalidPolicy") {
					return nil, nil
				}
				return nil, err
			}
			vulnerabilityAssessmentScanRecords = append(vulnerabilityAssessmentScanRecords, result.Value...)
		}
	}

	return vulnerabilityAssessmentScanRecords, nil
}

//// TRANSFORM FUNCTION

func idToServerName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem != nil {
		item := d.HydrateItem.(armsql.Database)
		return strings.Split(*item.ID, "/")[8], nil
	}
	return nil, nil
}
