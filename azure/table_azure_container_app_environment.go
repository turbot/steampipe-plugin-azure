package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v4"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureContainerAppEnvironment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_container_app_environment",
		Description: "Azure Container App Managed Environment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getContainerAppEnvironment,
			Tags: map[string]string{
				"service": "Microsoft.App",
				"action":  "managedEnvironments/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listContainerAppEnvironments,
			Tags: map[string]string{
				"service": "Microsoft.App",
				"action":  "managedEnvironments/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the managed environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the managed environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the managed environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The geo-location where the managed environment lives.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Kind of the environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning state of the environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "default_domain",
				Description: "Default domain name for the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DefaultDomain"),
			},
			{
				Name:        "static_ip",
				Description: "Static IP of the environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.StaticIP"),
			},
			{
				Name:        "deployment_errors",
				Description: "Any errors that occurred during deployment or deployment validation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DeploymentErrors"),
			},
			{
				Name:        "event_stream_endpoint",
				Description: "The endpoint of the eventstream of the environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.EventStreamEndpoint"),
			},
			{
				Name:        "infrastructure_resource_group",
				Description: "Name of the platform-managed resource group created for the managed environment to host infrastructure resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.InfrastructureResourceGroup"),
			},
			{
				Name:        "public_network_access",
				Description: "Property to allow or block all public traffic. Possible values are 'Enabled' and 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PublicNetworkAccess"),
			},
			{
				Name:        "zone_redundant",
				Description: "True if this managed environment is zone-redundant.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.ZoneRedundant"),
			},
			{
				Name:        "dapr_ai_connection_string",
				Description: "Application Insights connection string used by Dapr to export service to service communication telemetry.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DaprAIConnectionString"),
			},
			{
				Name:        "dapr_ai_instrumentation_key",
				Description: "Azure Monitor instrumentation key used by Dapr to export service to service communication telemetry.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DaprAIInstrumentationKey"),
			},
			{
				Name:        "app_logs_configuration",
				Description: "Cluster configuration which enables the log daemon to export app logs to a configured destination.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AppLogsConfiguration"),
			},
			{
				Name:        "custom_domain_configuration",
				Description: "Custom domain configuration for the environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.CustomDomainConfiguration"),
			},
			{
				Name:        "dapr_configuration",
				Description: "The configuration of the Dapr component.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.DaprConfiguration"),
			},
			{
				Name:        "keda_configuration",
				Description: "The configuration of the Keda component.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.KedaConfiguration"),
			},
			{
				Name:        "ingress_configuration",
				Description: "Ingress configuration for the managed environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.IngressConfiguration"),
			},
			{
				Name:        "peer_authentication",
				Description: "Peer authentication settings for the managed environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PeerAuthentication"),
			},
			{
				Name:        "peer_traffic_configuration",
				Description: "Peer traffic settings for the managed environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PeerTrafficConfiguration"),
			},
			{
				Name:        "vnet_configuration",
				Description: "Vnet configuration for the environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.VnetConfiguration"),
			},
			{
				Name:        "workload_profiles",
				Description: "Workload profiles configured for the managed environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.WorkloadProfiles"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "Private endpoint connections to the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PrivateEndpointConnections"),
			},
			{
				Name:        "identity",
				Description: "Managed identities for the managed environment to interact with other Azure services.",
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

func listContainerAppEnvironments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app_environment.listContainerAppEnvironments", "session_error", err)
		return nil, err
	}

	client, err := armappcontainers.NewManagedEnvironmentsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app_environment.listContainerAppEnvironments", "client_error", err)
		return nil, err
	}

	pager := client.NewListBySubscriptionPager(&armappcontainers.ManagedEnvironmentsClientListBySubscriptionOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_container_app_environment.listContainerAppEnvironments", "api_error", err)
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

func getContainerAppEnvironment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Empty check
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app_environment.getContainerAppEnvironment", "session_error", err)
		return nil, err
	}

	client, err := armappcontainers.NewManagedEnvironmentsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app_environment.getContainerAppEnvironment", "client_error", err)
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, &armappcontainers.ManagedEnvironmentsClientGetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_app_environment.getContainerAppEnvironment", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return &op.ManagedEnvironment, nil
	}

	return nil, nil
}
