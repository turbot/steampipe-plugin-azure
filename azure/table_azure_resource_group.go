package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureResourceGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_resource_group",
		Description: "Azure Resource Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ItemFromKey:       resourceGroupNameFromKey,
			Hydrate:           getResourceGroup,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listResourceGroups,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the resource group",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a resource group uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "Current state of the resource group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "managed_by",
				Description: "Contains ID of the resource that manages this resource group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Type of the resource group",
				Type:        proto.ColumnType_STRING,
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
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// ITEM FROM KEY

func resourceGroupNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &resources.Group{
		Name: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listResourceGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	resourcesClient := resources.NewGroupsClient(subscriptionID)
	resourcesClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := resourcesClient.List(ctx, "", nil)
		if err != nil {
			return nil, err
		}

		for _, resourceGroup := range result.Values() {
			d.StreamListItem(ctx, resourceGroup)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getResourceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	resourceGroup := h.Item.(*resources.Group)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	resourceGroupClient := resources.NewGroupsClient(subscriptionID)
	resourceGroupClient.Authorizer = session.Authorizer

	op, err := resourceGroupClient.Get(ctx, *resourceGroup.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
