package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureNetworkProfile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_network_profile",
		Description: "Azure Network Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getNetworkProfile,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkProfiles,
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
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the network profile resource. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfilePropertiesFormat.ProvisioningState"),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the network profile resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfilePropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "container_network_interfaces",
				Description: "List of child container network interfaces.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProfilePropertiesFormat.ContainerNetworkInterfaces"),
			},
			{
				Name:        "container_network_interface_configurations",
				Description: "List of child container network interface configurations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProfilePropertiesFormat.ContainerNetworkInterfaceConfigurations"),
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

//// LIST FUNCTION

func listNetworkProfiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := network.NewProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	result, err := client.ListAll(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listNetworkProfiles", "list", err)
		return nil, err
	}

	for _, profile := range result.Values() {
		d.StreamListItem(ctx, profile)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listNetworkProfiles", "list_paging", err)
			return nil, err
		}
		for _, profile := range result.Values() {
			d.StreamListItem(ctx, profile)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getNetworkProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getNetworkProfile")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := network.NewProfilesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	profile, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("getNetworkProfile", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if profile.ID != nil {
		return profile, nil
	}

	return nil, nil
}
