package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v4"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureContainerApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_container_app",
		Description: "Azure Container App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getContainerApp,
			Tags: map[string]string{
				"service": "Microsoft.App",
				"action":  "containerApps/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listContainerApps,
			Tags: map[string]string{
				"service": "Microsoft.App",
				"action":  "containerApps/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the container app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the container app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the container app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The geo-location where the container app lives.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Metadata to represent the container app kind, representing if a container app is workflowapp or functionapp.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_by",
				Description: "The fully qualified resource ID of the resource that manages this container app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the container app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "running_status",
				Description: "Running status of the container app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RunningStatus"),
			},
			{
				Name:        "environment_id",
				Description: "Resource ID of the environment the container app runs in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EnvironmentID"),
			},
			{
				Name:        "managed_environment_id",
				Description: "Deprecated. Resource ID of the container app's environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ManagedEnvironmentID"),
			},
			{
				Name:        "workload_profile_name",
				Description: "Workload profile name to pin for container app execution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.WorkloadProfileName"),
			},
			{
				Name:        "latest_revision_name",
				Description: "Name of the latest revision of the container app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.LatestRevisionName"),
			},
			{
				Name:        "latest_ready_revision_name",
				Description: "Name of the latest ready revision of the container app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.LatestReadyRevisionName"),
			},
			{
				Name:        "latest_revision_fqdn",
				Description: "Fully qualified domain name of the latest revision of the container app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.LatestRevisionFqdn"),
			},
			{
				Name:        "custom_domain_verification_id",
				Description: "Id used to verify domain name ownership.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.CustomDomainVerificationID"),
			},
			{
				Name:        "event_stream_endpoint",
				Description: "The endpoint of the eventstream of the container app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EventStreamEndpoint"),
			},
			{
				Name:        "active_revisions_mode",
				Description: "Controls how active revisions are handled for the container app: 'Multiple' allows several active revisions, 'Single' only one at a time.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Configuration.ActiveRevisionsMode"),
			},
			{
				Name:        "ingress_fqdn",
				Description: "Hostname the container app's ingress is reachable on.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Configuration.Ingress.Fqdn"),
			},
			{
				Name:        "ingress_external",
				Description: "True if the container app exposes an external HTTP endpoint.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.Configuration.Ingress.External"),
			},
			{
				Name:        "ingress_target_port",
				Description: "Target port in containers for traffic from ingress.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Configuration.Ingress.TargetPort"),
			},
			{
				Name:        "ingress_transport",
				Description: "Ingress transport protocol.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Configuration.Ingress.Transport"),
			},
			{
				Name:        "min_replicas",
				Description: "Minimum number of container replicas.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Template.Scale.MinReplicas"),
			},
			{
				Name:        "max_replicas",
				Description: "Maximum number of container replicas. Defaults to 10 if not set.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Template.Scale.MaxReplicas"),
			},
			{
				Name:        "configuration",
				Description: "Non versioned container app configuration properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Configuration"),
			},
			{
				Name:        "template",
				Description: "Container app versioned application definition.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Template"),
			},
			{
				Name:        "outbound_ip_addresses",
				Description: "Outbound IP addresses for the container app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.OutboundIPAddresses"),
			},
			{
				Name:        "identity",
				Description: "Managed identities for the container app to interact with other Azure services.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "extended_location",
				Description: "The complex type of the extended location.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "system_data",
				Description: "Azure Resource Manager metadata containing createdBy and modifiedBy information.",
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

func listContainerApps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app.listContainerApps", "session_error", err)
		return nil, err
	}

	client, err := armappcontainers.NewContainerAppsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app.listContainerApps", "client_error", err)
		return nil, err
	}

	pager := client.NewListBySubscriptionPager(&armappcontainers.ContainerAppsClientListBySubscriptionOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_container_app.listContainerApps", "api_error", err)
			return nil, err
		}

		for _, v := range page.Value {
			d.StreamListItem(ctx, v)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getContainerApp(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Empty check
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app.getContainerApp", "session_error", err)
		return nil, err
	}

	client, err := armappcontainers.NewContainerAppsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app.getContainerApp", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, &armappcontainers.ContainerAppsClientGetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app.getContainerApp", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return &op.ContainerApp, nil
	}

	return nil, nil
}
