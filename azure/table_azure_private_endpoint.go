package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzurePrivateEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_private_endpoint",
		Description: "Azure Private Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getPrivateEndpoint,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listPrivateEndpoints,
			KeyColumns: plugin.OptionalColumns([]string{"resource_group"}),
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the private endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "id",
				Description: "The ID of the private endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the private endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the private endpoint resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrivateEndpointProperties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "custom_network_interface_name",
				Description: "The custom name of the network interface attached to the private endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrivateEndpointProperties.CustomNetworkInterfaceName"),
			},
			{
				Name:        "location",
				Description: "The location of the private endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "extended_location",
				Description: "The extended location of the private endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnet",
				Description: "The ID of the subnet from which the private IP will be allocated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrivateEndpointProperties.Subnet"),
			},
			{
				Name:        "network_interfaces",
				Description: "An array of references to the network interfaces created for this private endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointProperties.NetworkInterfaces"),
			},
			{
				Name:        "private_link_service_connections",
				Description: "A grouping of information about the connection to the remote resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointProperties.PrivateLinkServiceConnections"),
			},
			{
				Name:        "manual_private_link_service_connections",
				Description: "A grouping of information about the connection to the remote resource. Used when the network admin does not have access to approve connections to the remote resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointProperties.ManualPrivateLinkServiceConnections"),
			},
			{
				Name:        "custom_dns_configs",
				Description: "An array of custom DNS configurations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointProperties.CustomDNSConfigs"),
			},
			{
				Name:        "application_security_groups",
				Description: "Application security groups in which the private endpoint IP configuration is included.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointProperties.ApplicationSecurityGroups"),
			},
			{
				Name:        "ip_configurations",
				Description: "A list of IP configurations of the private endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateEndpointProperties.IPConfigurations"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: "Tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: "The Azure region where the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: "The resource group in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

func listPrivateEndpoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_endpoint.listPrivateEndpoints", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := network.NewPrivateEndpointsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	resourceGroupName := h.Item.(resources.Group).Name

	if d.EqualsQualString("resource_group") != "" && d.EqualsQualString("resource_group") != *resourceGroupName {
		return nil, nil
	}

	result, err := client.List(ctx, *resourceGroupName)
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_endpoint.listPrivateEndpoints", "api_error", err)
		return nil, err
	}

	for _, privateEndpoint := range result.Values() {
		d.StreamListItem(ctx, privateEndpoint)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_private_endpoint.listPrivateEndpoints", "api_error_paging", err)
			return nil, err
		}

		for _, privateEndpoint := range result.Values() {
			d.StreamListItem(ctx, privateEndpoint)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

func getPrivateEndpoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_endpoint.getPrivateEndpoint", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := network.NewPrivateEndpointsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_endpoint.getPrivateEndpoint", "api_error", err)
		return nil, err
	}

	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
