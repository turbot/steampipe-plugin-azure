package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureNetworkSecurityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_network_security_group",
		Description: "Azure Network Security Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getNetworkSecurityGroup,
			Tags: map[string]string{
				"service": "Microsoft.Network",
				"action":  "networkSecurityGroups/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkSecurityGroups,
			Tags: map[string]string{
				"service": "Microsoft.Network",
				"action":  "networkSecurityGroups/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the network security group.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a network security group uniquely.",
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
				Description: "The resource type of the network security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The resource type of the network security group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "resource_guid",
				Description: "The resource GUID property of the network security group resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.ResourceGUID"),
			},
			{
				Name:        "default_security_rules",
				Description: "A list of default security rules of network security group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.DefaultSecurityRules"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the network security group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkSecurityGroupDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "flow_logs",
				Description: "A collection of references to flow log resources.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.FlowLogs"),
			},
			{
				Name:        "network_interfaces",
				Description: "A collection of references to network interfaces.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.NetworkInterfaces"),
			},
			{
				Name:        "security_rules",
				Description: "A list of security rules of network security group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.SecurityRules"),
			},
			{
				Name:        "subnets",
				Description: "A collection of references to subnets.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupPropertiesFormat.Subnets"),
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

func listNetworkSecurityGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkSecurityGroupClient := network.NewSecurityGroupsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkSecurityGroupClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkSecurityGroupClient, d.Connection)

	result, err := networkSecurityGroupClient.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, networkSecurityGroup := range result.Values() {
		d.StreamListItem(ctx, networkSecurityGroup)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, networkSecurityGroup := range result.Values() {
			d.StreamListItem(ctx, networkSecurityGroup)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getNetworkSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getNetworkSecurityGroup")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkSecurityGroupClient := network.NewSecurityGroupsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	networkSecurityGroupClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkSecurityGroupClient, d.Connection)

	op, err := networkSecurityGroupClient.Get(ctx, resourceGroup, name, "")
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

func listNetworkSecurityGroupDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listNetworkSecurityGroupDiagnosticSettings")
	id := *h.Item.(network.SecurityGroup).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.List(ctx, id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["id"] = i.ID
		}
		if i.Name != nil {
			objectMap["name"] = i.Name
		}
		if i.Type != nil {
			objectMap["type"] = i.Type
		}
		if i.DiagnosticSettings != nil {
			objectMap["properties"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}
	return diagnosticSettings, nil
}
