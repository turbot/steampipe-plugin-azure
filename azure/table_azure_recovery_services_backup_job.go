package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v3"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureRecoveryServicesBackupJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_recovery_services_backup_job",
		Description: "Azure Recovery Services Backup Job",
		List: &plugin.ListConfig{
			ParentHydrate: listRecoveryServicesVaults,
			Hydrate:       listRecoveryServicesBackupJobs,
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
				Description: "The recovery vault name.",
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
				Name:        "etag",
				Description: "Optional ETag.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},

			// JSON fields
			{
				Name:        "properties",
				Description: "JobResource properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(backupJobProperties),
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

type JobInfo struct {
	VaultName  *string
	ETag       *string
	Location   *string
	Properties armrecoveryservicesbackup.JobClassification
	Tags       map[string]*string
	ID         *string
	Name       *string
	Type       *string
}

//// LIST FUNCTION ////

func listRecoveryServicesBackupJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	vault := h.Item.(recoveryservices.Vault)

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

	clientFactory, err := armrecoveryservicesbackup.NewBackupJobsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_recovery_services_backup_job.listRecoveryServicesBackupJobs", "client_error", err)
		return nil, nil
	}
	pager := clientFactory.NewListPager(*vault.Name, strings.Split(*vault.ID, "/")[4], &armrecoveryservicesbackup.BackupJobsClientListOptions{Filter: nil,
		SkipToken: nil,
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_recovery_services_backup_job.listRecoveryServicesBackupJobs", "api_error", err)
			return nil, nil
		}

		for _, v := range page.Value {
			d.StreamListItem(ctx, JobInfo{
				ETag:       v.ETag,
				Location:   v.Location,
				Properties: v.Properties,
				Tags:       v.Tags,
				ID:         v.ID,
				Name:       v.Name,
				Type:       v.Type,
				VaultName:  vault.Name,
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

//// TRANSFORM FUNCTION

func backupJobProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(JobInfo)

	output := make(map[string]interface{})

	if data.Properties != nil {
		if data.Properties.GetJob() != nil {
			if data.Properties.GetJob().ActivityID != nil {
				output["ActivityID"] = data.Properties.GetJob().ActivityID
			}
			if data.Properties.GetJob().BackupManagementType != nil {
				output["BackupManagementType"] = data.Properties.GetJob().BackupManagementType
			}
			if data.Properties.GetJob().JobType != nil {
				output["JobType"] = data.Properties.GetJob().JobType
			}
			if data.Properties.GetJob().EndTime != nil {
				output["EndTime"] = data.Properties.GetJob().EndTime
			}
			if data.Properties.GetJob().EntityFriendlyName != nil {
				output["EntityFriendlyName"] = data.Properties.GetJob().EntityFriendlyName
			}
			if data.Properties.GetJob().Operation != nil {
				output["Operation"] = data.Properties.GetJob().Operation
			}
			if data.Properties.GetJob().StartTime != nil {
				output["StartTime"] = data.Properties.GetJob().StartTime
			}
			if data.Properties.GetJob().Status != nil {
				output["Status"] = data.Properties.GetJob().Status
			}
		}
	}
	return output, nil
}
