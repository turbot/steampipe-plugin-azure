package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerinstance/mgmt/containerinstance"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureContainerGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_container_group",
		Description: "Azure Container Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getContainerGroup,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "Invalid input", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listContainerGropus,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
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
				Name:        "os_type",
				Description: "The operating system type required by the containers in the container group. Possible values include: 'Windows', 'Linux'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerGroupProperties.OsType").Transform(transform.ToString),
			},
			{
				Name:        "restart_policy",
				Description: "Restart policy for all containers within the container group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerGroupProperties.RestartPolicy").Transform(transform.ToString),
			},
			{
				Name:        "sku",
				Description: "The SKU for a container group. Possible values include: 'Standard', 'Dedicated'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerGroupProperties.Sku").Transform(transform.ToString),
			},
			{
				Name:        "containers",
				Description: "The containers within the container group.",
				Hydrate:     extractContainerProperties,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.Containers"),
			},
			{
				Name:        "diagnostics",
				Description: "The diagnostic information for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.Diagnostics"),
			},
			{
				Name:        "dns_config",
				Description: "The DNS config information for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.DNSConfig"),
			},
			{
				Name:        "encryption_properties",
				Description: "The encryption properties for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.EncryptionProperties"),
			},
			{
				Name:        "identity",
				Description: "The identity of the container group, if configured.",
				Type:        proto.ColumnType_JSON,
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
				Name:        "instance_view",
				Description: "The instance view of the container group. Only valid in response.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.InstanceView"),
			},
			{
				Name:        "ip_address",
				Description: "The IP address type of the container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.IPAddress"),
			},
			{
				Name:        "subnet_ids",
				Description: "The subnet resource IDs for a container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.SubnetIds"),
			},
			{
				Name:        "volumes",
				Description: "The list of volumes that can be mounted by containers in this container group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerGroupProperties.Volumes"),
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listContainerGropus(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listContainerGropus")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerinstance.NewContainerGroupsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, group := range result.Values() {
		d.StreamListItem(ctx, group)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, group := range result.Values() {
			d.StreamListItem(ctx, group)
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getContainerGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getContainerGroup")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerinstance.NewContainerGroupsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func extractContainerProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("extractContainerProperties")
	containers := h.Item.(containerinstance.ContainerGroup).ContainerGroupProperties
	
	var containterInstances []map[string]interface{}
	for _, container := range *containers.Containers {
		objectMap := make(map[string]interface{})

		if container.Name != nil {
			objectMap["name"] = container.Name
		}
		if container.ContainerProperties != nil {
			if container.ContainerProperties.Image != nil {
				objectMap["image"] = container.ContainerProperties.Image
			}
			if container.ContainerProperties.Command != nil {
				objectMap["command"] = container.ContainerProperties.Command
			}
			if container.ContainerProperties.Ports != nil {
				objectMap["ports"] = container.ContainerProperties.Ports
			}
			if container.ContainerProperties.EnvironmentVariables != nil {
				objectMap["environmentVariables"] = container.ContainerProperties.EnvironmentVariables
			}
			if container.ContainerProperties.InstanceView != nil {
				objectMap["instanceView"] = container.ContainerProperties.InstanceView
			}
			if container.ContainerProperties.Resources != nil {
				objectMap["resources"] = container.ContainerProperties.Resources
			}
			if container.ContainerProperties.VolumeMounts != nil {
				objectMap["volumeMounts"] = container.ContainerProperties.VolumeMounts
			}
			if container.ContainerProperties.LivenessProbe != nil {
				objectMap["livenessProbe"] = container.ContainerProperties.LivenessProbe
			}
			if container.ContainerProperties.ReadinessProbe != nil {
				objectMap["readinessProbe"] = container.ContainerProperties.ReadinessProbe
			}
		}

		containterInstances = append(containterInstances, objectMap)
	}
	return containterInstances, nil
}