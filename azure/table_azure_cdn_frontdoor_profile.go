package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cdn/mgmt/cdn"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureCDNFrontDoorProfile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cdn_frontdoor_profile",
		Description: "Azure CDN Front Door Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAzureCDNFrontDoorProfile,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureCDNFrontDoorProfiles,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the CDN front door profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the CDN front door profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "sku_name",
				Description: "Name of the pricing tier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "kind",
				Description: "Kind of the profile. Used by portal to differentiate traditional CDN profile and new AFD profile.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_state",
				Description: "Resource status of the CDN front door profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfileProperties.ResourceState").Transform(transform.ToString),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning status of the CDN front door profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfileProperties.ProvisioningState"),
			},
			{
				Name:        "front_door_id",
				Description: "The ID of the front door.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfileProperties.FrontDoorID"),
			},
			{
				Name:        "origin_response_timeout_seconds",
				Description: "Send and receive timeout on forwarding request to the origin. When timeout is reached, the request fails and returns.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ProfileProperties.OriginResponseTimeoutSeconds"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: "Tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: "The Azure region where the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: "The resource group in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listAzureCDNFrontDoorProfiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_profile.listAzureCDNFrontDoorProfiles", "session_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := cdn.NewProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer
	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_profile.listAzureCDNFrontDoorProfiles", "api_error", err)
		return nil, err
	}

	for _, profile := range result.Values() {
		d.StreamListItem(ctx, profile)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_cdn_frontdoor_profile.listAzureCDNFrontDoorProfiles", "paging_error", err)
			return nil, err
		}
		for _, profile := range result.Values() {
			d.StreamListItem(ctx, profile)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getAzureCDNFrontDoorProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Return nil if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_profile.getAzureCDNFrontDoorProfile", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := cdn.NewProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	profile, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_profile.getAzureCDNFrontDoorProfile", "api_error", err)
		return nil, err
	}

	if profile.ID != nil {
		return profile, nil
	}

	return nil, nil
}
