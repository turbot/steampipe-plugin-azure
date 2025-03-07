package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/privatedns/mgmt/privatedns"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzurePrivateDNSZone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_private_dns_zone",
		Description: "Azure Private DNS Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getPrivateDNSZone,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listPrivateDNSZones,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the Private DNS zone.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a Private DNS zone uniquely.",
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
				Description: "The resource type of the Private DNS zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_number_of_record_sets",
				Description: "The maximum number of record sets that can be created in this Private DNS zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PrivateZoneProperties.MaxNumberOfRecordSets"),
			},
			{
				Name:        "number_of_record_sets",
				Description: "The current number of record sets in this Private DNS zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PrivateZoneProperties.NumberOfRecordSets").Transform(transform.ToString),
			},
			{
				Name:        "max_number_of_virtual_network_links",
				Description: "The maximum number of virtual networks that can be linked to this Private DNS zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PrivateZoneProperties.MaxNumberOfVirtualNetworkLinks"),
			},
			{
				Name:        "number_of_virtual_network_links",
				Description: "The current number of virtual networks that are linked to this Private DNS zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PrivateZoneProperties.NumberOfVirtualNetworkLinks"),
			},
			{
				Name:        "max_number_of_virtual_network_links_with_registration",
				Description: "The maximum number of virtual networks that can be linked to this Private DNS zone with registration enabled.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PrivateZoneProperties.MaxNumberOfVirtualNetworkLinksWithRegistration"),
			},
			{
				Name:        "number_of_virtual_network_links_with_registration",
				Description: "The current number of virtual networks that are linked to this Private DNS zone with registration enabled.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PrivateZoneProperties.NumberOfVirtualNetworkLinksWithRegistration"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the resource. Possible values include: `Creating`, `Updating`, `Deleting`, `Succeeded`, `Failed`, `Canceled`.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrivateZoneProperties.ProvisioningState"),
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

func listPrivateDNSZones(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_dns_zone.listPrivateDNSZones", "client_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	dnsClient := privatedns.NewPrivateZonesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	dnsClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &dnsClient, d.Connection)

	result, err := dnsClient.List(ctx, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_dns_zone.listPrivateDNSZones", "query_error", err)
		return nil, err
	}
	for _, dnsZone := range result.Values() {
		d.StreamListItem(ctx, dnsZone)
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

		for _, dnsZone := range result.Values() {
			d.StreamListItem(ctx, dnsZone)
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

func getPrivateDNSZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPrivateDNSZone")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_dns_zone.getPrivateDNSZone", "client_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	dnsClient := dns.NewZonesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	dnsClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &dnsClient, d.Connection)

	op, err := dnsClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_private_dns_zone.getPrivateDNSZone", "query_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
