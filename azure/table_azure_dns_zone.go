package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAzureDNSZone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_dns_zone",
		Description: "Azure DNS Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getDNSZone,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listDNSZones,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the DNS zone.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a DNS zone uniquely.",
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
				Description: "The resource type of the DNS zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_number_of_record_sets",
				Description: "The maximum number of record sets that can be created in this DNS zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ZoneProperties.MaxNumberOfRecordSets"),
			},
			{
				Name:        "max_number_of_records_per_record_set",
				Description: "The maximum number of records per record set that can be created in this DNS zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ZoneProperties.MaxNumberOfRecordsPerRecordSet"),
			},
			{
				Name:        "number_of_record_sets",
				Description: "The current number of record sets in this DNS zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ZoneProperties.NumberOfRecordSets").Transform(transform.ToString),
			},
			{
				Name:        "name_servers",
				Description: "The name servers for this DNS zone.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ZoneProperties.NameServers"),
			},
			{
				Name:        "zone_type",
				Description: "The type of this DNS zone (Public or Private).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ZoneProperties.ZoneType"),
			},
			{
				Name:        "registration_virtual_networks",
				Description: "A list of references to virtual networks that register hostnames in this DNS zone.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ZoneProperties.RegistrationVirtualNetworks"),
			},
			{
				Name:        "resolution_virtual_networks",
				Description: "A list of references to virtual networks that resolve records in this DNS zone.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ZoneProperties.ResolutionVirtualNetworks"),
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

func listDNSZones(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_dns_zone. listDNSZones", "client_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	dnsClient := dns.NewZonesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	dnsClient.Authorizer = session.Authorizer

	result, err := dnsClient.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, dnsZone := range result.Values() {
		d.StreamListItem(ctx, dnsZone)
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

		for _, dnsZone := range result.Values() {
			d.StreamListItem(ctx, dnsZone)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getDNSZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDNSZone")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	dnsClient := dns.NewZonesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	dnsClient.Authorizer = session.Authorizer

	op, err := dnsClient.Get(ctx, resourceGroup, name)
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
