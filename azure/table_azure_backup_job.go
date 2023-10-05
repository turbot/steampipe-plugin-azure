package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-01-01/backup"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureBackupJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_backup_job",
		Description: "Azure Backup Job",
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listAzureBackupJobs,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "404"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "vault_name",
					Require: plugin.Required,
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
				Transform:   transform.FromQual("vault_name"),
			},
			{
				Name:        "id",
				Description: "Resource Id represents the complete path to the resource.",
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
			{
				Name:        "time_created",
				Description: "The time when the disk was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DiskProperties.TimeCreated").Transform(convertDateToTime),
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

//// LIST FUNCTION ////

func listAzureBackupJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	resourceGroup := h.Item.(resources.Group)

	vaultName := d.EqualsQualString("vault_name")
	rgName := d.EqualsQualString("resource_group")

	if vaultName == "" {
		return nil, nil
	}

	if rgName != "" {
		if rgName != *resourceGroup.Name {
			return nil, nil
		}
	}

	subscriptionID := session.SubscriptionID
	client := backup.NewJobsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.List(ctx, vaultName, *resourceGroup.Name, "", "")
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, nil
		}
		return nil, err
	}

	for _, job := range result.Values() {
		d.StreamListItem(ctx, job)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotFound") {
				return nil, nil
			}
			return nil, err
		}

		for _, job := range result.Values() {
			d.StreamListItem(ctx, job)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
