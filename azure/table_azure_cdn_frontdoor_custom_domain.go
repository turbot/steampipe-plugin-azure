package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cdn/mgmt/cdn"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureCDNFrontDoorCustomDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_cdn_frontdoor_custom_domain",
		Description: "Azure CDN Front Door Custom Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "profile_name", "resource_group"}),
			Hydrate:    getAzureCDNFrontDoorCustomDomain,
			Tags: map[string]string{
				"service": "Microsoft.Cdn",
				"action":  "profiles/customDomains/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAzureCDNFrontDoorProfiles,
			Hydrate:       listAzureCDNFrontDoorCustomDomains,
			Tags: map[string]string{
				"service": "Microsoft.Cdn",
				"action":  "profiles/customDomains/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the custom domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "profile_name",
				Description: "The name of the CDN front door profile this custom domain belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractCDNProfileNameFromID),
			},
			{
				Name:        "host_name",
				Description: "The host name of the domain. Must be a domain name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.HostName"),
			},
			{
				Name:        "provisioning_state",
				Description: "Provisioning status of the custom domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "deployment_status",
				Description: "Deployment status of the custom domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DeploymentStatus"),
			},
			{
				Name:        "domain_validation_state",
				Description: "Provisioning substate showing the progress of the custom HTTPS enabling/disabling process.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DomainValidationState"),
			},
			{
				Name:        "azure_dns_zone_id",
				Description: "Resource reference to the Azure DNS zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AzureDNSZone.ID"),
			},
			{
				Name:        "pre_validated_custom_domain_resource_id",
				Description: "Resource reference to the Azure resource where custom domain ownership was prevalidated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PreValidatedCustomDomainResourceID.ID"),
			},
			{
				Name:        "tls_settings",
				Description: "The JSON object that contains the properties to secure the domain.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.TLSSettings"),
			},
			{
				Name:        "validation_properties",
				Description: "Values the customer needs to validate domain ownership.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ValidationProperties"),
			},
			{
				Name:        "extended_properties",
				Description: "Key-value pair representing migration properties for domains.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ExtendedProperties"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
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

func listAzureCDNFrontDoorCustomDomains(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	profile, ok := h.Item.(cdn.Profile)
	if !ok {
		return nil, nil
	}

	if profile.ID == nil || profile.Name == nil {
		return nil, nil
	}

	// Only process Azure Front Door profiles; skip classic CDN profiles
	if profile.Sku == nil {
		return nil, nil
	}
	skuName := string(profile.Sku.Name)
	if skuName != "Standard_AzureFrontDoor" && skuName != "Premium_AzureFrontDoor" {
		return nil, nil
	}

	// Extract resource group from profile ID
	splitID := strings.Split(*profile.ID, "/")
	if len(splitID) < 5 {
		return nil, nil
	}
	resourceGroup := splitID[4]
	profileName := *profile.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_custom_domain.listAzureCDNFrontDoorCustomDomains", "session_error", err)
		return nil, err
	}

	client, err := armcdn.NewAFDCustomDomainsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_custom_domain.listAzureCDNFrontDoorCustomDomains", "client_error", err)
		return nil, err
	}

	pager := client.NewListByProfilePager(resourceGroup, profileName, nil)
	for pager.More() {
		d.WaitForListRateLimit(ctx)

		page, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_cdn_frontdoor_custom_domain.listAzureCDNFrontDoorCustomDomains", "api_error", err)
			return nil, err
		}

		for _, domain := range page.Value {
			d.StreamListItem(ctx, domain)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getAzureCDNFrontDoorCustomDomain(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	profileName := d.EqualsQuals["profile_name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	if name == "" || profileName == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_custom_domain.getAzureCDNFrontDoorCustomDomain", "session_error", err)
		return nil, err
	}

	client, err := armcdn.NewAFDCustomDomainsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_custom_domain.getAzureCDNFrontDoorCustomDomain", "client_error", err)
		return nil, err
	}

	resp, err := client.Get(ctx, resourceGroup, profileName, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_cdn_frontdoor_custom_domain.getAzureCDNFrontDoorCustomDomain", "api_error", err)
		return nil, err
	}

	return &resp.AFDDomain, nil
}
