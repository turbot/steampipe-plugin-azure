package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerinstance/mgmt/containerinstance"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureContainerGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_container_group",
		Description: "Azure Container Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getContainerGroup,
			Tags: map[string]string{
				"service": "Microsoft.ContainerInstance",
				"action":  "containerGroups/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listContainerGroups,
			Tags: map[string]string{
				"service": "Microsoft.ContainerInstance",
				"action":  "containerGroups/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the container group. This only appears in the response.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerGroupProperties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "restart_policy",
				Description: "Restart policy for all containers within the container group. Possible values include: 'ContainerGroupRestartPolicyAlways', 'ContainerGroupRestartPolicyOnFailure', 'ContainerGroupRestartPolicyNever'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerGroupProperties.RestartPolicy"),
			},

			{
				Name:        "sku",
				Description: "The SKU for a container group. Possible values include: 'ContainerGroupSkuStandard', 'ContainerGroupSkuDedicated'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerGroupProperties.Sku"),
			},
			{
				Name:        "os_type",
				Description: "The operating system type required by the containers in the container group. Possible values include: 'OperatingSystemTypesWindows', 'OperatingSystemTypesLinux'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerGroupProperties.OsType"),
			},
			{
				Name:        "encryption_properties",
				Description: "The encryption settings of container registry.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.EncryptionProperties"),
			},
			{
				Name:        "containers",
				Description: "The containers within the container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.Containers"),
			},
			{
				Name:        "ip_address",
				Description: "The IP address type of the container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.IPAddress"),
			},
			{
				Name:        "volumes",
				Description: "The instance view of the container group. Only valid in response.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.Volumes"),
			},
			{
				Name:        "instance_view",
				Description: "The instance view of the container group. Only valid in response.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.InstanceView"),
			},
			{
				Name:        "diagnostics",
				Description: "The diagnostic information for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.Diagnostics"),
			},
			{
				Name:        "subnet_ids",
				Description: "The subnet resource IDs for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.SubnetIds"),
			},
			{
				Name:        "dns_config",
				Description: "The DNS config information for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.DNSConfig"),
			},
			{
				Name:        "init_containers",
				Description: "The init containers for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.InitContainers"),
			},
			{
				Name:        "image_registry_credentials",
				Description: "The image registry credentials by which the container group is created from.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.ImageRegistryCredentials"),
			},
			{
				Name:        "identity",
				Description: "The identity of the container group.",
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

func listContainerGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_group.listContainerGroups", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerinstance.NewContainerGroupsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_group.listContainerGroups", "api_error", err)
		return nil, err
	}

	for _, group := range result.Values() {
		d.StreamListItem(ctx, group)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, group := range result.Values() {
			d.StreamListItem(ctx, group)
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

func getContainerGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQualString("name")
	resourceGroup := d.EqualsQualString("resource_group")

	// Return nil, if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_group.getContainerGroup", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerinstance.NewContainerGroupsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_container_group.getContainerGroup", "api_error", err)
		return nil, err
	}

	return op, nil
}
