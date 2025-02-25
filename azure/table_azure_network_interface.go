package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureNetworkInterface(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_network_interface",
		Description: "Azure Network Interface",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getNetworkInterface,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkInterfaces,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a network interface uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "Providsioning state of the network interface resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "enable_accelerated_networking",
				Description: "Indicates whether the network interface is accelerated networking enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.EnableAcceleratedNetworking"),
			},
			{
				Name:        "enable_ip_forwarding",
				Description: "Indicates whether IP forwarding is enabled on this network interface.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.EnableIPForwarding"),
			},
			{
				Name:        "internal_dns_name_label",
				Description: "Relative DNS name for this NIC used for internal communications between VMs in the same virtual network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.InternalDNSNameLabel"),
			},
			{
				Name:        "internal_domain_name_suffix",
				Description: "Contains domain name suffix for the network interface.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.InternalDomainNameSuffix"),
			},
			{
				Name:        "internal_fqdn",
				Description: "Fully qualified DNS name supporting internal communications between VMs in the same virtual network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.InternalFqdn"),
			},
			{
				Name:        "is_primary",
				Description: "Indicates whether this is a primary network interface on a virtual machine.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.Primary"),
			},
			{
				Name:        "mac_address",
				Description: "The MAC address of the network interface.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.MacAddress"),
			},
			{
				Name:        "network_security_group_id",
				Description: "The reference to the NetworkSecurityGroup resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.NetworkSecurityGroup.ID"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the network interface resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "virtual_machine_id",
				Description: "The reference to a virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.VirtualMachine.ID"),
			},
			{
				Name:        "vnet_encryption_supported",
				Description: "Whether the virtual machine this NIC is attached to supports encryption.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.VnetEncryptionSupported"),
			},
			{
				Name:        "disable_tcp_state_tracking",
				Description: "Indicates whether to disable TCP state tracking.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("InterfacePropertiesFormat.DisableTCPStateTracking"),
			},
			{
				Name:        "workload_type",
				Description: "Workload type of the network interface for BareMetal resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.WorkloadType"),
			},
			{
				Name:        "nic_type",
				Description: "Type of network interface resource (e.g., Standard, Elastic).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.NicType"),
			},
			{
				Name:        "migration_phase",
				Description: "Migration phase of network interface resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.MigrationPhase"),
			},
			{
				Name:        "auxiliary_mode",
				Description: "Auxiliary mode of network interface resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InterfacePropertiesFormat.AuxiliaryMode"),
			},
			{
				Name:        "applied_dns_servers",
				Description: "A list of applied dns servers.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.AppliedDNSServers"),
			},
			{
				Name:        "dns_servers",
				Description: "A collection of DNS servers IP addresses.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.DNSSettings.DNSServers"),
			},
			{
				Name:        "hosted_workloads",
				Description: "A collection of references to linked BareMetal resources.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.HostedWorkloads"),
			},
			{
				Name:        "ip_configurations",
				Description: "A list of IPConfigurations of the network interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.IPConfigurations"),
			},
			{
				Name:        "tap_configurations",
				Description: "A collection of TapConfigurations of the network interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.TapConfigurations"),
			},
			{
				Name:        "dscp_configuration",
				Description: "A reference to the DSCP configuration to which the network interface is linked.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.DscpConfiguration"),
			},
			{
				Name:        "private_link_service",
				Description: "Private link service of the network interface resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.PrivateLinkService"),
			},
			{
				Name:        "private_endpoint",
				Description: "A reference to the private endpoint to which the network interface is linked.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InterfacePropertiesFormat.PrivateEndpoint"),
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

func listNetworkInterfaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewInterfacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkClient, d.Connection)

	result, err := networkClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, networkInterface := range result.Values() {
		d.StreamListItem(ctx, networkInterface)
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

		for _, networkInterface := range result.Values() {
			d.StreamListItem(ctx, networkInterface)
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

func getNetworkInterface(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getNetworkInterface")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewInterfacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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
