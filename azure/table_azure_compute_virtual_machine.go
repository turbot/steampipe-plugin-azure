package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeVirtualMachine(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_virtual_machine",
		Description: "Azure Compute Virtual Machine",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAzureComputeVirtualMachine,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureComputeVirtualMachines,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
			{
				Func:    getNicPublicIPs,
				Depends: []plugin.HydrateFunc{getVMNics},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the virtual machine",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "power_state",
				Description: "Specifies the power state of the vm",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAzureComputeVirtualMachineStatuses,
				Transform:   transform.FromField("Statuses").Transform(getPowerState),
			},

			{
				Name:        "private_ips",
				Description: "An array of private ip addesses associated with the vm.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVMNics,
				Transform:   transform.FromValue().Transform(getPrivateIpsFromIpconfig),
			},
			{
				Name:        "public_ips",
				Description: "An array of public ip addesses associated with the vm.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNicPublicIPs,
				Transform:   transform.FromValue(),
			},

			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource in Azure",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The virtual machine provisioning state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.ProvisioningState"),
			},
			{
				Name:        "vm_id",
				Description: "Specifies an unique ID for VM, which is a 128-bits identifier that is encoded and stored in all Azure IaaS VMs SMBIOS and can be read using platform BIOS commands",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.VMID"),
			},
			{
				Name:        "size",
				Description: "Specifies the size of the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.HardwareProfile.VMSize").Transform(transform.ToString),
			},
			{
				Name:        "admin_user_name",
				Description: "Specifies the name of the administrator account",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.AdminUsername"),
			},
			{
				Name:        "allow_extension_operations",
				Description: "Specifies whether extension operations should be allowed on the virtual machine",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.AllowExtensionOperations"),
			},
			{
				Name:        "availability_set_id",
				Description: "Specifies the ID of the availability set",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.AvailabilitySet.ID"),
			},
			{
				Name:        "billing_profile_max_price",
				Description: "Specifies the maximum price you are willing to pay for a Azure Spot VM/VMSS",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("VirtualMachineProperties.BillingProfile.MaxPrice"),
			},
			{
				Name:        "boot_diagnostics_enabled",
				Description: "Specifies whether boot diagnostics should be enabled on the Virtual Machine, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics.Enabled"),
			},
			{
				Name:        "boot_diagnostics_storage_uri",
				Description: "Contains the Uri of the storage account to use for placing the console output and screenshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics.StorageURI"),
			},
			{
				Name:        "computer_name",
				Description: "Specifies the host OS name of the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.ComputerName"),
			},
			{
				Name:        "disable_password_authentication",
				Description: "Specifies whether password authentication should be disabled",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.LinuxConfiguration.DisablePasswordAuthentication"),
			},
			{
				Name:        "eviction_policy",
				Description: "Specifies the eviction policy for the Azure Spot virtual machine and Azure Spot scale set",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.EvictionPolicy").Transform(transform.ToString),
			},
			{
				Name:        "image_exact_version",
				Description: "Specifies in decimal numbers, the version of platform image or marketplace image used to create the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.ExactVersion"),
			},
			{
				Name:        "image_id",
				Description: "Specifies the ID of the image to use",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.ID"),
			},
			{
				Name:        "image_offer",
				Description: "Specifies the offer of the platform image or marketplace image used to create the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Offer"),
			},
			{
				Name:        "image_publisher",
				Description: "Specifies the publisher of the image to use",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Publisher"),
			},
			{
				Name:        "image_sku",
				Description: "Specifies the sku of the image to use",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Sku"),
			},
			{
				Name:        "image_version",
				Description: "Specifies the version of the platform image or marketplace image used to create the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Version"),
			},
			{
				Name:        "managed_disk_id",
				Description: "Specifies the ID of the managed disk used by the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.ManagedDisk.ID"),
			},
			{
				Name:        "os_disk_caching",
				Description: "Specifies the caching requirements of the operating system disk used by the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.Caching").Transform(transform.ToString),
			},
			{
				Name:        "os_disk_create_option",
				Description: "Specifies how the virtual machine should be created",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.CreateOption").Transform(transform.ToString),
			},
			{
				Name:        "os_disk_name",
				Description: "Specifies the name of the operating system disk used by the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.Name"),
			},
			{
				Name:        "os_type",
				Description: "Specifies the type of the OS that is included in the disk if creating a VM from user-image or a specialized VHD",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.OsType").Transform(transform.ToString),
			},
			{
				Name:        "priority",
				Description: "Specifies the priority for the virtual machine",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.Priority").Transform(transform.ToString),
			},
			{
				Name:        "provision_vm_agent",
				Description: "Specifies whether virtual machine agent should be provisioned on the virtual machine",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.LinuxConfiguration.ProvisionVMAgent"),
			},
			{
				Name:        "require_guest_provision_signal",
				Description: "Specifies whether the guest provision signal is required to infer provision success of the virtual machine",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.RequireGuestProvisionSignal"),
			},
			{
				Name:        "ultra_ssd_enabled",
				Description: "Specifies whether managed disks with storage account type UltraSSD_LRS can be added to a virtual machine or virtual machine scale set, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.AdditionalCapabilities.UltraSSDEnabled"),
			},
			{
				Name:        "data_disks",
				Description: "A list of parameters that are used to add a data disk to a virtual machine",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.DataDisks"),
			},
			{
				Name:        "linux_configuration_ssh_public_keys",
				Description: "A list of ssh key configuration for a Linux OS",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.LinuxConfiguration.SSH.PublicKeys"),
			},
			{
				Name:        "network_interfaces",
				Description: "A list of resource Ids for the network interfaces associated with the virtual machine",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.NetworkProfile.NetworkInterfaces"),
			},
			{
				Name:        "secrets",
				Description: "A list of certificates that should be installed onto the virtual machine",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.Secrets"),
			},
			{
				Name:        "statuses",
				Description: "Specifies the resource status information",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAzureComputeVirtualMachineStatuses,
			},
			{
				Name:        "extensions",
				Description: "Specifies the details of VM Extensions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAzureComputeVirtualMachineExtensions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "zones",
				Description: "A list of virtual machine zones",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
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

			// Standard azure columns
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

//// LIST FUNCTION ////

func listAzureComputeVirtualMachines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeVirtualMachines")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachinesClient(subscriptionID)
	client.Authorizer = session.Authorizer
	pagesLeft := true

	for pagesLeft {
		result, err := client.ListAll(context.Background(), "")
		if err != nil {
			return nil, err
		}

		for _, virtualMachine := range result.Values() {
			d.StreamListItem(ctx, virtualMachine)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, nil
}

//// HYDRATE FUNCTION ////

func getAzureComputeVirtualMachine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeVirtualMachine")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachinesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(context.Background(), resourceGroup, name, "")
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

func getAzureComputeVirtualMachineStatuses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeVirtualMachineStatuses")

	virtualMachine := h.Item.(compute.VirtualMachine)
	resourceGroupName := strings.Split(string(*virtualMachine.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachinesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.InstanceView(context.Background(), resourceGroupName, *virtualMachine.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getVMNics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVMNics")

	vm := h.Item.(compute.VirtualMachine)
	var ipConfigs []network.InterfaceIPConfiguration

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	networkClient := network.NewInterfacesClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	for _, nicRef := range *vm.NetworkProfile.NetworkInterfaces {
		pathParts := strings.Split(string(*nicRef.ID), "/")
		resourceGroupName := pathParts[4]
		nicName := pathParts[len(pathParts)-1]

		nic, err := networkClient.Get(ctx, resourceGroupName, nicName, "")
		if err != nil {
			return nil, err
		}

		for _, ipConfig := range *nic.IPConfigurations {
			ipConfigs = append(ipConfigs, ipConfig)
		}
	}

	return ipConfigs, nil
}

func getNicPublicIPs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getNicPublicIPs")

	ipConfigs := h.HydrateResults["getVMNics"].([]network.InterfaceIPConfiguration)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	var publicIPs []string
	for _, ipConfig := range ipConfigs {
		if ipConfig.PublicIPAddress != nil && ipConfig.PublicIPAddress.ID != nil {
			publicIP, err := getNicPublicIP(ctx, session, *ipConfig.PublicIPAddress.ID)

			if err != nil {
				return nil, err
			}
			if publicIP.IPAddress != nil {
				publicIPs = append(publicIPs, *publicIP.IPAddress)
			}
		}
	}

	return publicIPs, nil
}

func getNicPublicIP(ctx context.Context, session *Session, id string) (network.PublicIPAddress, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getNicPublicIPs")

	pathParts := strings.Split(id, "/")
	resourceGroup := pathParts[4]
	name := pathParts[len(pathParts)-1]

	subscriptionID := session.SubscriptionID
	networkClient := network.NewPublicIPAddressesClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	return networkClient.Get(ctx, resourceGroup, name, "")
}

func getAzureComputeVirtualMachineExtensions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeVirtualMachineExtensions")

	virtualMachine := h.Item.(compute.VirtualMachine)
	resourceGroupName := strings.Split(string(*virtualMachine.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachineExtensionsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(context.Background(), resourceGroupName, *virtualMachine.Name, "")
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives the contents of VirtualMachineExtensionsListResult
	var extensions []map[string]interface{}

	for _, extension := range *op.Value {
		extensionMap := make(map[string]interface{})
		extensionMap["Id"] = extension.ID
		extensionMap["Name"] = extension.Name
		extensionMap["Type"] = extension.Type
		extensionMap["Location"] = extension.Location
		extensionMap["Publisher"] = extension.VirtualMachineExtensionProperties.Publisher
		extensionMap["TypeHandlerVersion"] = extension.VirtualMachineExtensionProperties.TypeHandlerVersion
		extensionMap["AutoUpgradeMinorVersion"] = extension.VirtualMachineExtensionProperties.AutoUpgradeMinorVersion
		extensionMap["EnableAutomaticUpgrade"] = extension.VirtualMachineExtensionProperties.EnableAutomaticUpgrade
		extensionMap["ForceUpdateTag"] = extension.VirtualMachineExtensionProperties.ForceUpdateTag
		extensionMap["Settings"] = extension.VirtualMachineExtensionProperties.Settings
		extensionMap["ProtectedSettings"] = extension.VirtualMachineExtensionProperties.ProtectedSettings
		extensionMap["ProvisioningState"] = extension.VirtualMachineExtensionProperties.ProvisioningState
		extensionMap["InstanceView"] = extension.VirtualMachineExtensionProperties.InstanceView
		extensionMap["Tags"] = extension.Tags
		extensions = append(extensions, extensionMap)
	}
	return extensions, nil
}

// TRANSFORM FUNCTIONS

func getPowerState(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPowerState", "d.Value", d.Value)

	if d.Value == nil {
		return nil, nil
	}
	statuses, ok := d.Value.(*[]compute.InstanceViewStatus)
	if ok != true {
		return nil, fmt.Errorf("Conversion failed for virtual machine statuses")
	}

	return getStatusFromCode(statuses, "PowerState"), nil
}

func getPrivateIpsFromIpconfig(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPrivateIpsFromIpconfig", "d.Value", d.Value)
	if d.Value == nil {
		return nil, nil
	}

	var ips []string
	ipConfigs, ok := d.Value.([]network.InterfaceIPConfiguration)
	if ok != true {
		return nil, fmt.Errorf("Conversion failed for virtual machine ip configs")
	}
	for _, ipConfig := range ipConfigs {
		ips = append(ips, *ipConfig.PrivateIPAddress)
	}

	return ips, nil
}

// UTILITY FUNCTIONS
func getStatusFromCode(statuses *[]compute.InstanceViewStatus, codeType string) string {
	for _, status := range *statuses {
		statusCode := types.SafeString(status.Code)

		if strings.HasPrefix(statusCode, codeType+"/") {
			return strings.SplitN(statusCode, "/", 2)[1]
		}
	}
	return ""
}
