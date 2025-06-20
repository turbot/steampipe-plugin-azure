package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureVirtualNetwork(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_virtual_network",
		Description: "Azure Virtual Network",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getVirtualNetwork,
			Tags: map[string]string{
				"service": "Microsoft.Network",
				"action":  "virtualNetworks/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listVirtualNetworks,
			Tags: map[string]string{
				"service": "Microsoft.Network",
				"action":  "virtualNetworks/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the virtual network",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a virtual network uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_ddos_protection",
				Description: "Indicates if DDoS protection is enabled for all the protected resources in the virtual network",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkPropertiesFormat.EnableDdosProtection"),
			},
			{
				Name:        "enable_vm_protection",
				Description: "Indicates if VM protection is enabled for all the subnets in the virtual network",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualNetworkPropertiesFormat.EnableVMProtection"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the virtual network resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "resource_guid",
				Description: "The resourceGuid property of the Virtual Network resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualNetworkPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "address_prefixes",
				Description: "A list of address blocks reserved for this virtual network in CIDR notation",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes"),
			},
			{
				Name:        "network_peerings",
				Description: "A list of peerings in a Virtual Network",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkPropertiesFormat.VirtualNetworkPeerings"),
			},
			{
				Name:        "subnets",
				Description: "A list of subnets in a Virtual Network",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualNetworkPropertiesFormat.Subnets"),
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

//// FETCH FUNCTIONS ////

func listVirtualNetworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	networkClient := network.NewVirtualNetworksClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkClient, d.Connection)

	result, err := networkClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, network := range result.Values() {
		d.StreamListItem(ctx, network)
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

		for _, network := range result.Values() {
			d.StreamListItem(ctx, network)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getVirtualNetwork(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVirtualNetwork")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewVirtualNetworksClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkClient, d.Connection)

	op, err := networkClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
