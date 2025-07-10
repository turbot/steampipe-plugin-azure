package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureDataProtectionBackupJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_protection_backup_job",
		Description: "Azure Data Protection Backup Job",
		List: &plugin.ListConfig{
			ParentHydrate: listAzureDataProtectionBackupVaults,
			Hydrate:       listAzureDataProtectionBackupJobs,
			Tags: map[string]string{
				"service": "Microsoft.DataProtection",
				"action":  "backupJobs/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "404"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "vault_name",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_group",
					Require: plugin.Optional,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Resource name associated with the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_name",
				Description: "The data protection vault name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Resource ID represents the complete path to the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/...",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activity_id",
				Description: "Job Activity Id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ActivityID"),
			},
			{
				Name:        "backup_instance_friendly_name",
				Description: "Name of the Backup Instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.BackupInstanceFriendlyName"),
			},
			{
				Name:        "data_source_id",
				Description: "ARM ID of the DataSource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DataSourceID"),
			},
			{
				Name:        "data_source_location",
				Description: "Location of the DataSource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DataSourceLocation"),
			},
			{
				Name:        "data_source_name",
				Description: "User Friendly Name of the DataSource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DataSourceName"),
			},
			{
				Name:        "data_source_type",
				Description: "Type of DataSource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DataSourceType"),
			},
			{
				Name:        "is_user_triggered",
				Description: "Indicates whether the job is adhoc(true) or scheduled(false).",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.IsUserTriggered"),
			},
			{
				Name:        "operation",
				Description: "Type of Job i.e. Backup:full/log/diff ;Restore:ALR/OLR; Tiering:Backup/Archive ; Management:ConfigureProtection/UnConfigure.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Operation"),
			},
			{
				Name:        "operation_category",
				Description: "Indicates the type of Job i.e. Backup/Restore/Tiering/Management.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.OperationCategory"),
			},
			{
				Name:        "progress_enabled",
				Description: "Indicates whether progress is enabled for the job.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.ProgressEnabled"),
			},
			{
				Name:        "source_resource_group",
				Description: "Resource Group Name of the Datasource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SourceResourceGroup"),
			},
			{
				Name:        "source_subscription_id",
				Description: "SubscriptionId corresponding to the DataSource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SourceSubscriptionID"),
			},
			{
				Name:        "start_time",
				Description: "StartTime of the job (in UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.StartTime"),
			},
			{
				Name:        "status",
				Description: "Status of the job like InProgress/Success/Failed/Cancelled/SuccessWithWarning.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Status"),
			},
			{
				Name:        "data_source_set_name",
				Description: "Data Source Set Name of the DataSource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DataSourceSetName"),
			},
			{
				Name:        "destination_data_store_name",
				Description: "Destination Data Store Name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DestinationDataStoreName"),
			},
			{
				Name:        "duration",
				Description: "Total run time of the job. ISO 8601 format.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Duration"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Etag"),
			},
			{
				Name:        "source_data_store_name",
				Description: "Source Data Store Name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SourceDataStoreName"),
			},
			{
				Name:        "backup_instance_id",
				Description: "ARM ID of the Backup Instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.BackupInstanceID"),
			},
			{
				Name:        "end_time",
				Description: "EndTime of the job (in UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.EndTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "policy_id",
				Description: "ARM ID of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PolicyID"),
			},
			{
				Name:        "policy_name",
				Description: "Name of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PolicyName"),
			},
			{
				Name:        "progress_url",
				Description: "Url which contains job's progress.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProgressURL"),
			},
			{
				Name:        "restore_type",
				Description: "Indicates the sub type of operation i.e. in case of Restore it can be ALR/OLR.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RestoreType"),
			},
			{
				Name:        "error_details",
				Description: "A List, detailing the errors related to the job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ErrorDetails"),
			},
			{
				Name:        "supported_actions",
				Description: "List of supported actions.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.SupportedActions"),
			},
			{
				Name:        "extended_info",
				Description: "Extended Information about the job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ExtendedInfo"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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

type DataProtectionJobInfo struct {
	VaultName  *string
	Location   *string
	Properties *armdataprotection.AzureBackupJob
	ID         *string
	Name       *string
	Type       *string
}

//// LIST FUNCTION ////

func listAzureDataProtectionBackupJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_data_protection_backup_job.listAzureDataProtectionBackupJobs", "session_error", err)
		return nil, err
	}
	vault := h.Item.(*armdataprotection.BackupVaultResource)

	vaultName := d.EqualsQualString("vault_name")
	rgName := d.EqualsQualString("resource_group")

	if vaultName != "" {
		if vaultName != *vault.Name {
			return nil, nil
		}
	}

	if rgName != "" {
		if rgName != strings.Split(*vault.ID, "/")[4] {
			return nil, nil
		}
	}

	clientFactory, err := armdataprotection.NewJobsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_data_protection_backup_job.listAzureDataProtectionBackupJobs", "client_error", err)
		return nil, nil
	}
	pager := clientFactory.NewListPager(strings.Split(*vault.ID, "/")[4], *vault.Name, &armdataprotection.JobsClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_data_protection_backup_job.listAzureDataProtectionBackupJobs", "api_error", err)
			return nil, nil
		}

		for _, v := range page.Value {
			d.StreamListItem(ctx, DataProtectionJobInfo{
				Properties: v.Properties,
				ID:         v.ID,
				Name:       v.Name,
				Type:       v.Type,
				VaultName:  vault.Name,
				Location:   vault.Location,
			})

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
