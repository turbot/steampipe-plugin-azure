package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAKSVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_kubernetes_service_version",
		Description: "Azure Kubernetes Service Version",
		List: &plugin.ListConfig{
			Hydrate: listAKSVersions,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "location",
					Require: plugin.Required,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "The major.minor version of Kubernetes release.",
			},
			{
				Name:        "is_preview",
				Description: "Whether Kubernetes version is currently in preview.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "capabilities",
				Description: "Capabilities on this Kubernetes version.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "patch_versions",
				Description: "Patch versions of Kubernetes release.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Version"),
			},

			// Azure standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("location"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAKSVersions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	location := d.EqualsQualString("location")

	// Empty Check
	if location == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_kubernetes_service_version.listAKSVersions", "session_error", err)
		return nil, err
	}

	client, err := armcontainerservice.NewManagedClustersClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	result, err := client.ListKubernetesVersions(ctx, location, nil)

	for _, op := range result.Values {
		d.StreamListItem(ctx, op)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}
