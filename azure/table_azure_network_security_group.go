package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureNetworkSecurityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_network_security_group",
		Description: "Azure Network Security Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getNetworkSecurityGroup,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkSecurityGroups,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the network security group",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a network security group uniquely",
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
				Description: "The resource type of the network security group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The resource type of the network security group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the network security group resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "default_security_rules",
				Description: "A list of default security rules of network security group",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.DefaultSecurityRules"),
			},
			{
				Name:        "flow_logs",
				Description: "A collection of references to flow log resources",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.FlowLogs"),
			},
			{
				Name:        "network_interfaces",
				Description: "A collection of references to network interfaces",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.NetworkInterfaces"),
			},
			{
				Name:        "security_rules",
				Description: "A list of security rules of network security group",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.SecurityRules"),
			},
			{
				Name:        "subnets",
				Description: "A collection of references to subnets",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.Subnets"),
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

//// FETCH FUNCTIONS ////

func listNetworkSecurityGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	NetworkSecurityGroupClient := network.NewSecurityGroupsClient(subscriptionID)
	NetworkSecurityGroupClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := NetworkSecurityGroupClient.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		for _, networkSecurityGroup := range result.Values() {
			d.StreamListItem(ctx, networkSecurityGroup)
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getNetworkSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getNetworkSecurityGroup")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	NetworkSecurityGroupClient := network.NewSecurityGroupsClient(subscriptionID)
	NetworkSecurityGroupClient.Authorizer = session.Authorizer

	op, err := NetworkSecurityGroupClient.Get(ctx, resourceGroup, name, "")
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
