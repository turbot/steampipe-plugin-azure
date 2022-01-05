package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureComputeVirtualMachineScaleSetVm(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_virtual_machine_scale_set_vm",
		Description: "Azure Compute Virtual Machine Scale Set VM",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"scale_set_name", "resource_group", "instance_id"}),
			Hydrate:           getAzureComputeVirtualMachineScaleSetVm,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAzureComputeVirtualMachineScaleSets,
			Hydrate:       listAzureComputeVirtualMachineScaleSetVms,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the scale set VM.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scale_set_name",
				Description: "Name of the scale set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "instance_id",
				Description: "The virtual machine instance ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "latest_model_applied",
				Description: "Specifies whether the latest model has been applied to the virtual machine.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.LatestModelApplied"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource in Azure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_type",
				Description: "Specifies that the image or disk that is being used was licensed on-premises.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.LicenseType"),
			},
			{
				Name:        "location",
				Description: "The location of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "model_definition_applied",
				Description: "Specifies whether the model applied to the virtual machine is the model of the virtual machine scale set or the customized model for the virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.ModelDefinitionApplied"),
			},
			{
				Name:        "sku_name",
				Description: "The sku name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_capacity",
				Description: "Specifies the capacity of virtual machines in a scale set virtual machine.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_tier",
				Description: "Specifies the tier of virtual machines in a scale set virtual machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "vm_id",
				Description: "Azure virtual machine unique ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.VMID"),
			},
			{
				Name:        "additional_capabilities",
				Description: "Specifies additional capabilities enabled or disabled on the virtual machine in the scale set. For instance: whether the virtual machine has the capability to support attaching managed data disks with UltraSSD_LRS storage account type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.AdditionalCapabilities"),
			},
			{
				Name:        "availability_set",
				Description: "Specifies information about the availability set that the virtual machine should be assigned to. Virtual machines specified in the same availability set are allocated to different nodes to maximize availability.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.AvailabilitySet"),
			},
			{
				Name:        "plan",
				Description: "Specifies information about the marketplace image used to create the virtual machine.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "protection_policy",
				Description: "Specifies the protection policy of the virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.ProtectionPolicy"),
			},
			{
				Name:        "resources",
				Description: "The virtual machine child extension resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "Resource tags.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "upgrade_policy",
				Description: "The upgrade policy for the scale set.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.UpgradePolicy"),
			},
			{
				Name:        "virtual_machine_diagnostics_profile",
				Description: "Specifies the boot diagnostic settings state.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.DiagnosticsProfile"),
			},
			{
				Name:        "virtual_machine_hardware_profile",
				Description: "Specifies the hardware settings for the virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.HardwareProfile"),
			},
			{
				Name:        "virtual_machine_network_profile",
				Description: "Specifies properties of the network interfaces of the virtual machines.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.NetworkProfile.NetworkInterfaces"),
			},
			{
				Name:        "virtual_machine_network_profile_configuration",
				Description: "Specifies the network profile configuration of the virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.NetworkProfileConfiguration"),
			},
			{
				Name:        "virtual_machine_os_profile",
				Description: "Specifies the operating system settings for the virtual machines.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.OsProfile"),
			},
			{
				Name:        "virtual_machine_security_profile",
				Description: "Specifies the Security related profile settings for the virtual machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.SecurityProfile"),
			},
			{
				Name:        "virtual_machine_storage_profile",
				Description: "SSpecifies the storage settings for the virtual machine disks.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetVMProperties.StorageProfile"),
			},
			{
				Name:        "zones",
				Description: "The Logical zone list for scale set.",
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
		},
	}
}

type ScaleSetVMInfo struct {
	ScaleSetName string
	compute.VirtualMachineScaleSetVM
}

//// LIST FUNCTION

func listAzureComputeVirtualMachineScaleSetVms(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureComputeVirtualMachineScaleSetVms")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	scaleSet := h.Item.(compute.VirtualMachineScaleSet)
	resourceGroupName := strings.ToLower(strings.Split(*scaleSet.ID, "/")[4])

	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(context.Background(), resourceGroupName, *scaleSet.Name, "", "", "")
	if err != nil {
		plugin.Logger(ctx).Error("Error", "listAzureComputeVirtualMachineScaleSetVms", err)
		return nil, err
	}

	for _, scaleSetVm := range result.Values() {
		d.StreamListItem(ctx, ScaleSetVMInfo{*scaleSet.Name, scaleSetVm})
	}
	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, scaleSetVm := range result.Values() {
			d.StreamListItem(ctx, ScaleSetVMInfo{*scaleSet.Name, scaleSetVm})
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTION

func getAzureComputeVirtualMachineScaleSetVm(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeVirtualMachineScaleSetVm")

	scaleSetName := d.KeyColumnQuals["scale_set_name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	instanceId := d.KeyColumnQuals["instance_id"].GetStringValue()

	if scaleSetName == "" || resourceGroup == "" || instanceId == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("Error", "getAzureComputeVirtualMachineScaleSetVm", "connection", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(context.Background(), resourceGroup, scaleSetName, instanceId, "")
	if err != nil {
		plugin.Logger(ctx).Error("Error", "getAzureComputeVirtualMachineScaleSetVm", err)
		return nil, err
	}

	if op.ID != nil {
		return ScaleSetVMInfo{scaleSetName, op}, nil
	}

	return nil, nil
}
