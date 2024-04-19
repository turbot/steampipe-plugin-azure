package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzurePublicIP(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_public_ip",
		Description: "Azure Public IP",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getPublicIP,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listPublicIPs,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the public ip",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a public ip uniquely",
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
				Description: "The resource type of the public ip",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The resource type of the public ip",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "ddos_custom_policy_id",
				Description: "The DDoS custom policy associated with the public IP",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.DdosSettings.DdosCustomPolicy.ID"),
			},
			{
				Name:        "ddos_settings_protected_ip",
				Description: "Indicates whether DDoS protection is enabled on the public IP, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.DdosSettings.ProtectedIP"),
			},
			{
				Name:        "ddos_settings_protection_coverage",
				Description: "The DDoS protection policy customizability of the public IP",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.DdosSettings.ProtectionCoverage").Transform(transform.ToString),
			},
			{
				Name:        "dns_settings_domain_name_label",
				Description: "Contains the domain name label",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel"),
			},
			{
				Name:        "dns_settings_fqdn",
				Description: "The Fully Qualified Domain Name of the A DNS record associated with the public IP",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.DNSSettings.Fqdn"),
			},
			{
				Name:        "dns_settings_reverse_fqdn",
				Description: "Contains the reverse FQDN",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.DNSSettings.ReverseFqdn"),
			},
			{
				Name:        "idle_timeout_in_minutes",
				Description: "The idle timeout of the public IP address",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes"),
			},
			{
				Name:        "ip_address",
				Description: "The IP address associated with the public IP address resource",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.IPAddress"),
			},
			{
				Name:        "ip_configuration_id",
				Description: "Contains the IP configuration ID",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.IPConfiguration.ID"),
			},
			{
				Name:        "public_ip_address_version",
				Description: "Contains the public IP address version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.PublicIPAddressVersion").Transform(transform.ToString),
			},
			{
				Name:        "public_ip_allocation_method",
				Description: "Contains the public IP address allocation method",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.PublicIPAllocationMethod").Transform(transform.ToString),
			},
			{
				Name:        "public_ip_prefix_id",
				Description: "The Public IP Prefix this Public IP Address should be allocated from",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.PublicIPPrefix.ID"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the public ip resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "sku_name",
				Description: "Name of a public IP address SKU",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "ip_tags",
				Description: "A list of tags associated with the public IP address",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PublicIPAddressPropertiesFormat.IPTags"),
			},
			{
				Name:        "zones",
				Description: "A collection of availability zones denoting the IP allocated for the resource needs to come from",
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
		}),
	}
}

//// FETCH FUNCTIONS ////

func listPublicIPs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	resourceGroup := h.Item.(resources.Group).Name

	networkClient := network.NewPublicIPAddressesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// ListAll API doesn't return any value so changed to List API
	result, err := networkClient.List(ctx, *resourceGroup)
	if err != nil {
		return nil, err
	}

	for _, publicIP := range result.Values() {
		d.StreamListItem(ctx, publicIP)
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

		for _, publicIP := range result.Values() {
			d.StreamListItem(ctx, publicIP)
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

func getPublicIP(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPublicIP")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewPublicIPAddressesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
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
