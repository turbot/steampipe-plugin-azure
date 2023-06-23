package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-09-01/containerservice"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAKSOrchestractor(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_kubernetes_service_version",
		Description: "Azure Kubernetes Service Version",
		List: &plugin.ListConfig{
			Hydrate: listAKSOrchestractors,
			// KeyColumns: plugin.AllColumns([]string{"region", "resource_type"}),
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "region",
					Require: plugin.Required,
				},
				{
					Name:    "resource_type",
					Require: plugin.Optional,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the orchestrator version profile list result.",
			},
			{
				Name:        "id",
				Description: "Id of the orchestrator version profile list result.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Type of the orchestrator version profile list result.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "orchestrator_type",
				Description: "The orchestrator type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "orchestrator_version",
				Description: "Orchestrator version (major, minor, patch).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default",
				Description: "Installed by default if version is not specified.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_preview",
				Description: "Whether Kubernetes version is currently in preview.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "resource_type",
				Description: "Whether Kubernetes version is currently in preview.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_type"),
			},
			{
				Name:        "upgrades",
				Description: "The list of available upgrade versions.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.FromQual("region"),
			},
		}),
	}
}

type OrchestratorInfo struct {
	// ID - READ-ONLY; Id of the orchestrator version profile list result.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Name of the orchestrator version profile list result.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Type of the orchestrator version profile list result.
	Type *string `json:"type,omitempty"`
	// OrchestratorType - Orchestrator type.
	OrchestratorType *string `json:"orchestratorType,omitempty"`
	// OrchestratorVersion - Orchestrator version (major, minor, patch).
	OrchestratorVersion *string `json:"orchestratorVersion,omitempty"`
	// Default - Installed by default if version is not specified.
	Default *bool `json:"default,omitempty"`
	// IsPreview - Whether Kubernetes version is currently in preview.
	IsPreview *bool `json:"isPreview,omitempty"`
	// Upgrades - The list of available upgrade versions.
	Upgrades *[]containerservice.OrchestratorProfile `json:"upgrades,omitempty"`
}

//// LIST FUNCTION

func listAKSOrchestractors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("region")
	resourceType := d.EqualsQualString("resource_type")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_kubernetes_service_version.listAKSOrchestractors", "session_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID

	containerserviceClient := containerservice.NewContainerServicesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	containerserviceClient.Authorizer = session.Authorizer

	result, err := containerserviceClient.ListOrchestrators(ctx, region, resourceType)
	if err != nil {
		plugin.Logger(ctx).Error("azure_kubernetes_service_version.listAKSOrchestractors", "api_error", err)
		return nil, err
	}

	for _, op := range *result.Orchestrators {
		d.StreamListItem(ctx, &OrchestratorInfo{
			ID:                  result.ID,
			Name:                result.Name,
			Type:                result.Type,
			OrchestratorType:    op.OrchestratorType,
			OrchestratorVersion: op.OrchestratorVersion,
			Default:             op.Default,
			IsPreview:           op.IsPreview,
			Upgrades:            op.Upgrades,
		})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}
