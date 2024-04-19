package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/backup"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/recoveryservices"

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
	Properties backup.BasicJob
	Tags       map[string]*string
	ID         *string
	Name       *string
	Type       *string
}

//// LIST FUNCTION ////

func listRecoveryServicesBackupJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	vault := h.Item.(recoveryservices.Vault)

	subscriptionID := session.SubscriptionID
	backupClient := backup.NewJobsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	backupClient.Authorizer = session.Authorizer
	result, err := backupClient.List(ctx, *vault.Name, strings.Split(*vault.ID, "/")[4], "", "")
	if err != nil {
		return nil, err
	}
	for _, vault := range result.Values() {
		d.StreamListItem(ctx, vault)
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
		for _, v := range result.Values() {
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
	return nil, err
}

//// TRANSFORM FUNCTION

func backupJobProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(JobInfo)

	output := make(map[string]interface{})

	if data.Properties != nil {
		job, flag := data.Properties.AsJob()
		if flag {
			if job.ActivityID != nil {
				output["ActivityID"] = job.ActivityID
			}
			output["BackupManagementType"] = job.BackupManagementType
			output["JobType"] = job.JobType
			if job.EndTime != nil {
				output["EndTime"] = job.EndTime
			}
			if job.EntityFriendlyName != nil {
				output["EntityFriendlyName"] = job.EntityFriendlyName
			}
			if job.Operation != nil {
				output["Operation"] = job.Operation
			}
			if job.StartTime != nil {
				output["StartTime"] = job.StartTime
			}
			if job.Status != nil {
				output["Status"] = job.Status
			}
		}
	}
	return output, nil
}
