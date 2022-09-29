package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/guestconfiguration/mgmt/2020-06-25/guestconfiguration"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeVirtualMachine(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_virtual_machine",
		Description: "Azure Compute Virtual Machine",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getComputeVirtualMachine,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeVirtualMachines,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:    getNicPublicIPs,
				Depends: []plugin.HydrateFunc{getVMNics},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the virtual machine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "power_state",
				Description: "Specifies the power state of the vm.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getComputeVirtualMachineInstanceView,
				Transform:   transform.FromField("Statuses").Transform(getPowerState),
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource in Azure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The virtual machine provisioning state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.ProvisioningState"),
			},
			{
				Name:        "vm_id",
				Description: "Specifies an unique ID for VM, which is a 128-bits identifier that is encoded and stored in all Azure IaaS VMs SMBIOS and can be read using platform BIOS commands.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.VMID"),
			},
			{
				Name:        "size",
				Description: "Specifies the size of the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.HardwareProfile.VMSize").Transform(transform.ToString),
			},
			{
				Name:        "admin_user_name",
				Description: "Specifies the name of the administrator account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.AdminUsername"),
			},
			{
				Name:        "allow_extension_operations",
				Description: "Specifies whether extension operations should be allowed on the virtual machine.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.AllowExtensionOperations"),
			},
			{
				Name:        "availability_set_id",
				Description: "Specifies the ID of the availability set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.AvailabilitySet.ID"),
			},
			{
				Name:        "billing_profile_max_price",
				Description: "Specifies the maximum price you are willing to pay for a Azure Spot VM/VMSS.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("VirtualMachineProperties.BillingProfile.MaxPrice"),
			},
			{
				Name:        "boot_diagnostics_enabled",
				Description: "Specifies whether boot diagnostics should be enabled on the Virtual Machine, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics.Enabled"),
			},
			{
				Name:        "boot_diagnostics_storage_uri",
				Description: "Contains the Uri of the storage account to use for placing the console output and screenshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics.StorageURI"),
			},
			{
				Name:        "computer_name",
				Description: "Specifies the host OS name of the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.ComputerName"),
			},
			{
				Name:        "disable_password_authentication",
				Description: "Specifies whether password authentication should be disabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.LinuxConfiguration.DisablePasswordAuthentication"),
			},
			{
				Name:        "enable_automatic_updates",
				Description: "Indicates whether automatic updates is enabled for the windows virtual machine.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.WindowsConfiguration.EnableAutomaticUpdates"),
			},
			{
				Name:        "eviction_policy",
				Description: "Specifies the eviction policy for the Azure Spot virtual machine and Azure Spot scale set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.EvictionPolicy").Transform(transform.ToString),
			},
			{
				Name:        "image_exact_version",
				Description: "Specifies in decimal numbers, the version of platform image or marketplace image used to create the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.ExactVersion"),
			},
			{
				Name:        "image_id",
				Description: "Specifies the ID of the image to use.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.ID"),
			},
			{
				Name:        "image_offer",
				Description: "Specifies the offer of the platform image or marketplace image used to create the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Offer"),
			},
			{
				Name:        "image_publisher",
				Description: "Specifies the publisher of the image to use.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Publisher"),
			},
			{
				Name:        "image_sku",
				Description: "Specifies the sku of the image to use.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Sku"),
			},
			{
				Name:        "image_version",
				Description: "Specifies the version of the platform image or marketplace image used to create the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.ImageReference.Version"),
			},
			{
				Name:        "managed_disk_id",
				Description: "Specifies the ID of the managed disk used by the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.ManagedDisk.ID"),
			},
			{
				Name:        "os_disk_caching",
				Description: "Specifies the caching requirements of the operating system disk used by the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.Caching").Transform(transform.ToString),
			},
			{
				Name:        "os_disk_create_option",
				Description: "Specifies how the virtual machine should be created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.CreateOption").Transform(transform.ToString),
			},
			{
				Name:        "os_disk_name",
				Description: "Specifies the name of the operating system disk used by the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.Name"),
			},
			{
				Name:        "os_disk_vhd_uri",
				Description: "Specifies the virtual hard disk's uri.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.Vhd.URI").Transform(transform.ToString),
			},
			{
				Name:        "os_name",
				Description: "The Operating System running on the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getComputeVirtualMachineInstanceView,
			},
			{
				Name:        "os_version",
				Description: "The version of Operating System running on the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getComputeVirtualMachineInstanceView,
			},
			{
				Name:        "os_type",
				Description: "Specifies the type of the OS that is included in the disk if creating a VM from user-image or a specialized VHD.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.StorageProfile.OsDisk.OsType").Transform(transform.ToString),
			},
			{
				Name:        "priority",
				Description: "Specifies the priority for the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.Priority").Transform(transform.ToString),
			},
			{
				Name:        "provision_vm_agent",
				Description: "Specifies whether virtual machine agent should be provisioned on the virtual machine for linux configuration.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.LinuxConfiguration.ProvisionVMAgent"),
			},
			{
				Name:        "provision_vm_agent_windows",
				Description: "Specifies whether virtual machine agent should be provisioned on the virtual machine for windows configuration.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.WindowsConfiguration.ProvisionVMAgent"),
			},
			{
				Name:        "require_guest_provision_signal",
				Description: "Specifies whether the guest provision signal is required to infer provision success of the virtual machine.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.RequireGuestProvisionSignal"),
			},
			{
				Name:        "ultra_ssd_enabled",
				Description: "Specifies whether managed disks with storage account type UltraSSD_LRS can be added to a virtual machine or virtual machine scale set, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineProperties.AdditionalCapabilities.UltraSSDEnabled"),
			},
			{
				Name:        "time_zone",
				Description: "Specifies the time zone of the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.WindowsConfiguration.TimeZone"),
			},
			{
				Name:        "additional_unattend_content",
				Description: "Specifies additional base-64 encoded XML formatted information that can be included in the Unattend.xml file, which is used by windows setup.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.WindowsConfiguration.AdditionalUnattendContent"),
			},
			{
				Name:        "data_disks",
				Description: "A list of parameters that are used to add a data disk to a virtual machine.",
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
				Description: "A list of resource Ids for the network interfaces associated with the virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.NetworkProfile.NetworkInterfaces"),
			},
			{
				Name:        "patch_settings",
				Description: "Specifies settings related to in-guest patching (KBs).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.WindowsConfiguration.PatchSettings"),
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
				Name:        "secrets",
				Description: "A list of certificates that should be installed onto the virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.Secrets"),
			},
			{
				Name:        "statuses",
				Description: "Specifies the resource status information.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeVirtualMachineInstanceView,
			},
			{
				Name:        "extensions",
				Description: "Specifies the details of VM Extensions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAzureComputeVirtualMachineExtensions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "guest_configuration_assignments",
				Description: "Guest configuration assignments for a virtual machine.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listComputeVirtualMachineGuestConfigurationAssignments,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity",
				Description: "The identity of the virtual machine, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_profile",
				Description: "Specifies the security related profile settings for the virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.SecurityProfile"),
			},
			{
				Name:        "win_rm",
				Description: "Specifies the windows remote management listeners. This enables remote windows powershell.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineProperties.OsProfile.WindowsConfiguration.WinRM"),
			},
			{
				Name:        "zones",
				Description: "A list of virtual machine zones.",
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

//// LIST FUNCTION ////

func listComputeVirtualMachines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeVirtualMachines")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachinesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.ListAll(ctx, "")
	if err != nil {
		return nil, err
	}

	for _, virtualMachine := range result.Values() {
		d.StreamListItem(ctx, virtualMachine)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, virtualMachine := range result.Values() {
			d.StreamListItem(ctx, virtualMachine)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION ////

func getComputeVirtualMachine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeVirtualMachine")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachinesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name, "")
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

func getComputeVirtualMachineInstanceView(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeVirtualMachineInstanceView")

	virtualMachine := h.Item.(compute.VirtualMachine)
	resourceGroupName := strings.Split(string(*virtualMachine.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachinesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.InstanceView(ctx, resourceGroupName, *virtualMachine.Name)
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
	networkClient := network.NewInterfacesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	for _, nicRef := range *vm.NetworkProfile.NetworkInterfaces {
		pathParts := strings.Split(string(*nicRef.ID), "/")
		resourceGroupName := pathParts[4]
		nicName := pathParts[len(pathParts)-1]

		nic, err := networkClient.Get(ctx, resourceGroupName, nicName, "")
		if err != nil {
			return nil, err
		}

		ipConfigs = append(ipConfigs, *nic.IPConfigurations...)
	}

	return ipConfigs, nil
}

func getNicPublicIPs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getNicPublicIPs")

	// Interface IP Configuration will be nil if getVMNics returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getVMNics"] == nil {
		return nil, nil
	}
	
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
	networkClient := network.NewPublicIPAddressesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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
	client := compute.NewVirtualMachineExtensionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, resourceGroupName, *virtualMachine.Name, "")
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
		extensionMap["ExtensionType"] = extension.VirtualMachineExtensionProperties.Type
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

func listComputeVirtualMachineGuestConfigurationAssignments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeVirtualMachineGuestConfigurationAssignments")

	virtualMachine := h.Item.(compute.VirtualMachine)
	resourceGroupName := strings.Split(string(*virtualMachine.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := guestconfiguration.NewAssignmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// SDK does not support pagination yet
	op, err := client.List(ctx, resourceGroupName, *virtualMachine.Name)
	if err != nil {
		// API throws 404 error if vm does not have any guest configuration assignments
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listComputeVirtualMachineGuestConfigurationAssignments", "get", err)
		return nil, err
	}

	var assignments []map[string]interface{}

	// If we return the API response directly, the output will not provide all the data for Guest Configuration Assignment
	for _, configAssignment := range *op.Value {
		objectMap := make(map[string]interface{})
		if configAssignment.ID != nil {
			objectMap["id"] = configAssignment.ID
		}
		if configAssignment.Name != nil {
			objectMap["name"] = configAssignment.Name
		}
		if configAssignment.Location != nil {
			objectMap["location"] = configAssignment.Location
		}
		if configAssignment.Type != nil {
			objectMap["type"] = configAssignment.Type
		}
		if configAssignment.Properties != nil {
			if configAssignment.Properties.TargetResourceID != nil {
				objectMap["targetResourceID"] = configAssignment.Properties.TargetResourceID
			}
			if configAssignment.Properties.TargetResourceID != nil {
				objectMap["lastComplianceStatusChecked"] = configAssignment.Properties.LastComplianceStatusChecked
			}
			if configAssignment.Properties.ComplianceStatus != "" {
				objectMap["complianceStatus"] = configAssignment.Properties.ComplianceStatus
			}
			if configAssignment.Properties.LatestReportID != nil {
				objectMap["latestReportID"] = configAssignment.Properties.LatestReportID
			}
			if configAssignment.Properties.Context != nil {
				objectMap["context"] = configAssignment.Properties.Context
			}
			if configAssignment.Properties.AssignmentHash != nil {
				objectMap["assignmentHash"] = configAssignment.Properties.AssignmentHash
			}
			if configAssignment.Properties.ProvisioningState != "" {
				objectMap["provisioningState"] = configAssignment.Properties.ProvisioningState
			}
			if configAssignment.Properties.GuestConfiguration != nil {
				objectMap["guestConfiguration"] = configAssignment.Properties.GuestConfiguration
			}
			if configAssignment.Properties.LatestAssignmentReport != nil {
				objectMap["latestAssignmentReport"] = configAssignment.Properties.LatestAssignmentReport
			}
		}
		assignments = append(assignments, objectMap)
	}

	return assignments, nil
}

// TRANSFORM FUNCTIONS

func getPowerState(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPowerState", "d.Value", d.Value)

	if d.Value == nil {
		return nil, nil
	}
	statuses, ok := d.Value.(*[]compute.InstanceViewStatus)
	if !ok {
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
	if !ok {
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
