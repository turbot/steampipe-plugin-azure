package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureComputeVirtualMachineScaleSetNetworkInterface(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_virtual_machine_scale_set_network_interface",
		Description: "Azure Compute Virtual Machine Scale Set Network Interface",
		List: &plugin.ListConfig{
			ParentHydrate: listAzureComputeVirtualMachineScaleSets,
			Hydrate:       listAzureComputeVirtualMachineScaleSetInterfaces,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the scale set network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scale_set_name",
				Description: "Name of the scale set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractScaleSetFromID),
			},
			{
				Name:        "id",
				Description: "The unique ID identifying the resource in a subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the network interface resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractScaleSetNetworkInterfacePrperties, "ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource in Azure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mac_address",
				Description: "The MAC address of the network interface.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Interface.MacAddress"),
			},
			{
				Name:        "enable_accelerated_networking",
				Description: "If the network interface has accelerated networking enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "EnableAcceleratedNetworking"),
			},
			{
				Name:        "enable_ip_forwarding",
				Description: "Indicates whether IP forwarding is enabled on this network interface.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "EnableIPForwarding"),
			},
			{
				Name:        "primary",
				Description: "Whether this is a primary network interface on a virtual machine.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "Primary"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the network interface resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "ResourceGUID"),
			},
			{
				Name:        "dns_settings",
				Description: "The DNS settings in network interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "DNSSettings"),
			},
			{
				Name:        "hosted_workloads",
				Description: "A list of references to linked BareMetal resources.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "HostedWorkloads"),
			},
			{
				Name:        "ip_configurations",
				Description: "A list of IP configurations of the network interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "IPConfigurations"),
			},
			{
				Name:        "network_security_group",
				Description: "The reference to the NetworkSecurityGroup resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "NetworkSecurityGroup"),
			},
			{
				Name:        "private_endpoint",
				Description: "A reference to the private endpoint to which the network interface is linked.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "PrivateEndpoint"),
			},
			{
				Name:        "virtual_machine",
				Description: "The reference to a virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractScaleSetNetworkInterfaceProperties, "VirtualMachine"),
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

func listAzureComputeVirtualMachineScaleSetInterfaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_compute_virtual_machine_scale_set_network_interface.listAzureComputeVirtualMachineScaleSetInterfaces", "session_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := network.NewInterfacesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	scaleSetinfo := h.Item.(compute.VirtualMachineScaleSet)
	resourceGroupName := strings.Split(string(*scaleSetinfo.ID), "/")[4]

	result, err := client.ListVirtualMachineScaleSetNetworkInterfaces(ctx, resourceGroupName, *scaleSetinfo.Name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_compute_virtual_machine_scale_set_network_interface.listAzureComputeVirtualMachineScaleSetInterfaces", "api_error", err)
		return nil, err
	}

	for _, scaleSetNetworkInterfacce := range result.Values() {
		d.StreamListItem(ctx, scaleSetNetworkInterfacce)
	}
	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, scaleSetNetworkInterfacce := range result.Values() {
			d.StreamListItem(ctx, scaleSetNetworkInterfacce)
		}
	}
	return nil, nil
}

//// TRANSFORM FUNCTION

func extractScaleSetNetworkInterfaceProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	networkInterface := d.HydrateItem.(network.Interface)
	param := d.Param.(string)

	objectMap := make(map[string]interface{})

	if networkInterface.VirtualMachine != nil {
		objectMap["VirtualMachine"] = *networkInterface.VirtualMachine
	}
	if *networkInterface.ResourceGUID != "" {
		objectMap["ResourceGUID"] = networkInterface.ResourceGUID
	}
	if networkInterface.ProvisioningState != "" {
		objectMap["ProvisioningState"] = networkInterface.ProvisioningState
	}
	if networkInterface.NetworkSecurityGroup != nil {
		objectMap["NetworkSecurityGroup"] = networkInterface.NetworkSecurityGroup
	}
	if networkInterface.IPConfigurations != nil {
		objectMap["IPConfigurations"] = networkInterface.IPConfigurations
	}
	if networkInterface.TapConfigurations != nil {
		objectMap["TapConfigurations"] = networkInterface.TapConfigurations
	}
	if networkInterface.DNSSettings != nil {
		objectMap["DNSSettings"] = networkInterface.DNSSettings
	}
	if networkInterface.MacAddress != nil {
		objectMap["MacAddress"] = networkInterface.MacAddress
	}
	if networkInterface.Primary != nil {
		objectMap["Primary"] = networkInterface.Primary
	}
	if networkInterface.EnableAcceleratedNetworking != nil {
		objectMap["EnableAcceleratedNetworking"] = networkInterface.EnableAcceleratedNetworking
	}
	if networkInterface.HostedWorkloads != nil {
		objectMap["HostedWorkloads"] = networkInterface.HostedWorkloads
	}
	if networkInterface.HostedWorkloads != nil {
		objectMap["HostedWorkloads"] = networkInterface.HostedWorkloads
	}
	if networkInterface.ResourceGUID != nil {
		objectMap["ResourceGUID"] = networkInterface.ResourceGUID
	}
	if networkInterface.ProvisioningState != "" {
		objectMap["ProvisioningState"] = networkInterface.ProvisioningState
	}

	if val, ok := objectMap[param]; ok {
		return val, nil
	}
	return nil, nil
}

func extractScaleSetFromID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)

	// Common resource properties
	splitID := strings.Split(id, "/")
	scaleSetName := splitID[8]
	scaleSetName = strings.ToLower(scaleSetName)
	return scaleSetName, nil
}
