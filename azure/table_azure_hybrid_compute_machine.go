package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/hybridcompute/mgmt/hybridcompute"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAzureHybridComputeMachine(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_hybrid_compute_machine",
		Description: "Azure Hybrid Compute Machine",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getHybridComputeMachine,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listHybridComputeMachines,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "status",
				Description: "The status of the hybrid machine agent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.Status"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the hybrid machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ad_fqdn",
				Description: "Specifies the AD fully qualified display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.AdFqdn"),
			},
			{
				Name:        "agent_version",
				Description: "The hybrid machine agent full version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.AgentVersion"),
			},
			{
				Name:        "client_public_key",
				Description: "Public Key that the client provides to be used during initial resource onboarding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.ClientPublicKey"),
			},
			{
				Name:        "dns_fqdn",
				Description: "Specifies the DNS fully qualified display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.DNSFqdn"),
			},
			{
				Name:        "display_name",
				Description: "Specifies the hybrid machine display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.DisplayName"),
			},
			{
				Name:        "domain_name",
				Description: "Specifies the Windows domain name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.DomainName"),
			},
			{
				Name:        "last_status_change",
				Description: "The time of the last status change.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("MachinePropertiesModel.LastStatusChange").Transform(convertDateToTime),
			},
			{
				Name:        "machine_fqdn",
				Description: "Specifies the hybrid machine FQDN.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.MachineFqdn"),
			},
			{
				Name:        "os_name",
				Description: "The Operating System running on the hybrid machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.OsName"),
			},
			{
				Name:        "os_sku",
				Description: "Specifies the Operating System product SKU.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.OsSku"),
			},
			{
				Name:        "os_version",
				Description: "The version of Operating System running on the hybrid machine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.OsVersion"),
			},
			{
				Name:        "vm_id",
				Description: "Specifies the hybrid machine unique ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.VMID"),
			},
			{
				Name:        "vm_uuid",
				Description: "Specifies the Arc Machine's unique SMBIOS ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachinePropertiesModel.VMUUID"),
			},
			{
				Name:        "error_details",
				Description: "Details about the error state.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MachinePropertiesModel.ErrorDetails"),
			},
			{
				Name:        "extensions",
				Description: "The extensions of the compute machine.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listHybridComputeMachineExtensions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity",
				Description: "The identity of the compute machine.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "location_data",
				Description: "The metadata pertaining to the geographic location of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MachinePropertiesModel.LocationData"),
			},
			{
				Name:        "machine_properties_extensions",
				Description: "The machine properties extensions of the compute machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MachinePropertiesModel.Extensions"),
			},
			{
				Name:        "os_profile",
				Description: "Specifies the operating system settings for the hybrid machine.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MachinePropertiesModel.OsProfile"),
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
				Transform:   transform.FromField("Tags"),
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

func listHybridComputeMachines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := hybridcompute.NewMachinesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListBySubscription(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listHybridComputeMachines", "list", err)
		return nil, err
	}

	for _, machine := range result.Values() {
		d.StreamListItem(ctx, machine)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listHybridComputeMachines", "list_paging", err)
			return nil, err
		}
		for _, machine := range result.Values() {
			d.StreamListItem(ctx, machine)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getHybridComputeMachine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getHybridComputeMachine")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := hybridcompute.NewMachinesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	machine, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("getHybridComputeMachine", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if machine.ID != nil {
		return machine, nil
	}

	return nil, nil
}

func listHybridComputeMachineExtensions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	machine := h.Item.(hybridcompute.Machine)
	resourceGroup := strings.Split(*machine.ID, "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := hybridcompute.NewMachineExtensionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	var extensions []map[string]interface{}

	result, err := client.List(ctx, resourceGroup, *machine.Name, "")
	if err != nil {
		plugin.Logger(ctx).Error("listHybridComputeMachineExtensions", "list", err)
		return nil, err
	}

	for _, extension := range result.Values() {
		extensions = append(extensions, extractComputeMachineExtensions(extension))
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listHybridComputeMachineExtensions", "list_paging", err)
			return nil, err
		}
		for _, extension := range result.Values() {
			extensions = append(extensions, extractComputeMachineExtensions(extension))
		}
	}

	return extensions, nil
}

func extractComputeMachineExtensions(extension hybridcompute.MachineExtension) map[string]interface{} {
	objectMap := make(map[string]interface{})
	if extension.ID != nil {
		objectMap["id"] = extension.ID
	}
	if extension.Name != nil {
		objectMap["name"] = extension.Name
	}
	if extension.Type != nil {
		objectMap["type"] = extension.Type
	}
	if extension.ProvisioningState != nil {
		objectMap["provisioningState"] = extension.ProvisioningState
	}
	if extension.MachineExtensionProperties != nil {
		objectMap["machineExtensionProperties"] = extension.MachineExtensionProperties
	}
	return objectMap
}
