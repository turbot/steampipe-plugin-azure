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

func tableAzureVirtualMachineScaleSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_virtual_machine_scale_set",
		Description: "Azure Virtual Machine Scale Set",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAzureVirtualMachineScaleSet,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureVirtualMachineScaleSets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
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
				Name:        "type",
				Description: "The type of the resource in Azure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "single_placement_group",
				Description: "When true this limits the scale set to a single placement group, of max size 100 virtual machines.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.SinglePlacementGroup"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.ProvisioningState"),
			},
			{
				Name:        "overprovision",
				Description: "Specifies whether the Virtual Machine Scale Set should be overprovisioned.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.Overprovision"),
			},
			{
				Name:        "do_not_run_extensions_on_overprovisioned_vms",
				Description: "When Overprovision is enabled, extensions are launched only on the requested number of VMs which are finally kept.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.DoNotRunExtensionsOnOverprovisionedVMs"),
			},
			{
				Name:        "unique_id",
				Description: "Specifies the ID which uniquely identifies a Virtual Machine Scale Set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.UniqueID"),
			},
			{
				Name:        "platform_fault_domain_count",
				Description: "Fault Domain count for each placement group.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.PlatformFaultDomainCount"),
			},
			{
				Name:        "sku_name",
				Description: "The sku name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "Specifies the tier of virtual machines in a scale set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier"),
			},
			{
				Name:        "sku_capacity",
				Description: "Specifies the tier of virtual machines in a scale set.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "extensions",
				Description: "Specifies the details of VM Scale Set Extensions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAzureComputeVirtualMachineScalesetExtensions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "plan",
				Description: "Specifies information about the marketplace image used to create the virtual machine.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "zones",
				Description: "The Logical zone list for scale set.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "The identity of the virtual machine scale set, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "virtual_machine_os_profile",
				Description: "Specifies the operating system settings for the virtual machines in the scale set.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.VirtualMachineProfile.OsProfile"),
			},
			{
				Name:        "virtual_machine_storage_profile",
				Description: "Specifies the storage settings for the virtual machine disks.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile"),
			},
			{
				Name:        "virtual_machine_network_profile",
				Description: "Specifies properties of the network interfaces of the virtual machines in the scale set.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.VirtualMachineProfile.NetworkProfile"),
			},
			{
				Name:        "virtual_machine_security_profile",
				Description: "Specifies the Security related profile settings for the virtual machines in the scale set.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.VirtualMachineProfile.SecurityProfile"),
			},
			{
				Name:        "virtual_machine_diagnostics_profile",
				Description: "Specifies the boot diagnostic settings state.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.VirtualMachineProfile.DiagnosticsProfile"),
			},
			{
				Name:        "virtual_machine_extension_profile",
				Description: "Specifies a collection of settings for extensions installed on virtual machines in the scale set.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.VirtualMachineProfile.ExtensionProfile"),
			},
			{
				Name:        "upgrade_policy",
				Description: "The upgrade policy for the scale set.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.UpgradePolicy"),
			},
			{
				Name:        "scale_in_policy",
				Description: "Specifies the scale-in policy that decides which virtual machines are chosen for removal when a Virtual Machine Scale Set is scaled-in.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VirtualMachineScaleSetProperties.ScaleInPolicy"),
			},
			{
				Name:        "tags_src",
				Description: "Resource tags.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
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

//// LIST FUNCTION

func listAzureVirtualMachineScaleSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAzureVirtualMachineScaleSet")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachineScaleSetsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListAll(context.Background())
	if err != nil {
		return nil, err
	}

	for _, scaleSet := range result.Values() {
		d.StreamListItem(ctx, scaleSet)
	}
	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, scaleSet := range result.Values() {
			d.StreamListItem(ctx, scaleSet)
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTION

func getAzureVirtualMachineScaleSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureVirtualMachineScaleSet")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachineScaleSetsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(context.Background(), resourceGroup, name)
	if err != nil {
		return nil, err
	}

	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func getAzureComputeVirtualMachineScalesetExtensions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeVirtualMachineScalesetExtensions")

	virtualMachineScaleSet := h.Item.(compute.VirtualMachineScaleSet)
	resourceGroupName := strings.Split(string(*virtualMachineScaleSet.ID), "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewVirtualMachineScaleSetExtensionsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(context.Background(), resourceGroupName, *virtualMachineScaleSet.Name)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives the contents of VirtualMachineScaleSetExtensionsListResult
	var extensions []map[string]interface{}

	for _, extension := range op.Values() {
		extensionMap := make(map[string]interface{})
		extensionMap["Id"] = extension.ID
		extensionMap["Name"] = extension.Name
		extensionMap["Type"] = extension.Type
		extensionMap["ProvisionAfterExtensions"] = extension.ProvisionAfterExtensions
		extensionMap["Publisher"] = extension.VirtualMachineScaleSetExtensionProperties.Publisher
		extensionMap["ExtensionType"] = extension.VirtualMachineScaleSetExtensionProperties.Type
		extensionMap["TypeHandlerVersion"] = extension.VirtualMachineScaleSetExtensionProperties.TypeHandlerVersion
		extensionMap["AutoUpgradeMinorVersion"] = extension.VirtualMachineScaleSetExtensionProperties.AutoUpgradeMinorVersion
		extensionMap["EnableAutomaticUpgrade"] = extension.VirtualMachineScaleSetExtensionProperties.EnableAutomaticUpgrade
		extensionMap["ForceUpdateTag"] = extension.VirtualMachineScaleSetExtensionProperties.ForceUpdateTag
		extensionMap["Settings"] = extension.VirtualMachineScaleSetExtensionProperties.Settings
		extensionMap["ProtectedSettings"] = extension.VirtualMachineScaleSetExtensionProperties.ProtectedSettings
		extensionMap["ProvisioningState"] = extension.VirtualMachineScaleSetExtensionProperties.ProvisioningState
		plugin.Logger(ctx).Trace("Extensions ==>", extensionMap)
		extensions = append(extensions, extensionMap)
	}
	return extensions, nil
}