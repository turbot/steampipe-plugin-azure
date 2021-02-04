package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureNetworkInterface(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_network_interface",
		Description: "Azure Network Interface",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getNetworkInterface,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkInterfaces,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the network interface",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a network interface uniquely",
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
				Description: "The resource type of the network interface",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Providsioning state of the network interface resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "enable_accelerated_networking",
				Description: "Indicates whether the network interface is accelerated networking enabled",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.EnableAcceleratedNetworking"),
			},
			{
				Name:        "enable_ip_forwarding",
				Description: "Indicates whether IP forwarding is enabled on this network interface",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.EnableIPForwarding"),
			},
			{
				Name:        "internal_dns_name_label",
				Description: "Relative DNS name for this NIC used for internal communications between VMs in the same virtual network",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.InternalDNSNameLabel"),
			},
			{
				Name:        "internal_domain_name_suffix",
				Description: "Contains domain name suffix for the network interface",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.InternalDomainNameSuffix"),
			},
			{
				Name:        "internal_fqdn",
				Description: "Fully qualified DNS name supporting internal communications between VMs in the same virtual network",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.InternalFqdn"),
			},
			{
				Name:        "is_primary",
				Description: "Indicates whether this is a primary network interface on a virtual machine",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.Primary"),
			},
			{
				Name:        "mac_address",
				Description: "The MAC address of the network interface",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.MacAddress"),
			},
			{
				Name:        "network_security_group_id",
				Description: "The reference to the NetworkSecurityGroup resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.NetworkSecurityGroup.ID"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the network interface resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "virtual_machine_id",
				Description: "The reference to a virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.VirtualMachine.ID"),
			},
			{
				Name:        "applied_dns_servers",
				Description: "A list of applied dns servers",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.AppliedDNSServers"),
			},
			{
				Name:        "dns_servers",
				Description: "A collection of DNS servers IP addresses",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.DNSServers"),
			},
			{
				Name:        "hosted_workloads",
				Description: "A collection of references to linked BareMetal resources",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.HostedWorkloads"),
			},
			{
				Name:        "ip_configurations",
				Description: "A list of IPConfigurations of the network interface",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.IPConfigurations"),
			},
			{
				Name:        "tap_configurations",
				Description: "A collection of TapConfigurations of the network interface",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.TapConfigurations"),
			},

			// Standard columns
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

//// FETCH FUNCTIONS ////

func listNetworkInterfaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewInterfacesClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := networkClient.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		for _, networkInterface := range result.Values() {
			d.StreamListItem(ctx, networkInterface)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getNetworkInterface(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getNetworkInterface")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewInterfacesClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

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
