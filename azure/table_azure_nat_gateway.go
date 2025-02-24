package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureNatGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_nat_gateway",
		Description: "Azure NAT Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getNatGateway,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listNatGateways,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the nat gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a nat gateway uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "idle_timeout_in_minutes",
				Description: "The idle timeout of the nat gateway.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromP(extractNatGatewayProperties, "idleTimeoutInMinutes"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the nat gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractNatGatewayProperties, "provisioningState"),
			},
			{
				Name:        "resource_guid",
				Description: "The provisioning state of the nat gateway resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(extractNatGatewayProperties, "resourceGUID"),
			},
			{
				Name:        "sku_name",
				Description: "The nat gateway SKU.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "type",
				Description: "The resource type of the nat gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_ip_addresses",
				Description: "An array of public ip addresses associated with the nat gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractNatGatewayProperties, "publicIpAddresses"),
			},
			{
				Name:        "public_ip_prefixes",
				Description: "An array of public ip prefixes associated with the nat gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractNatGatewayProperties, "publicIpPrefixes"),
			},
			{
				Name:        "subnets",
				Description: "An array of references to the subnets using this nat gateway resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(extractNatGatewayProperties, "subnets"),
			},
			{
				Name:        "zones",
				Description: "A list of availability zones denoting the zone in which Nat Gateway should be deployed.",
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

func listNatGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_nat_gateway.listNatGateways", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewNatGatewaysClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkClient, d.Connection)

	result, err := networkClient.ListAll(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_nat_gateway.listNatGateways", "api_error", err)
		return nil, err
	}

	for _, natGateway := range result.Values() {
		d.StreamListItem(ctx, natGateway)
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

		for _, natGateWay := range result.Values() {
			d.StreamListItem(ctx, natGateWay)
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

func getNatGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_nat_gateway.getNatGateway", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	networkClient := network.NewNatGatewaysClient(subscriptionID)
	networkClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &networkClient, d.Connection)

	op, err := networkClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("azure_nat_gateway.getNatGateway", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func extractNatGatewayProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	gateway := d.HydrateItem.(network.NatGateway)
	param := d.Param.(string)

	objectMap := make(map[string]interface{})

	if gateway.IdleTimeoutInMinutes != nil {
		objectMap["idleTimeoutInMinutes"] = *gateway.IdleTimeoutInMinutes
	}
	if *gateway.ResourceGUID != "" {
		objectMap["resourceGUID"] = gateway.ResourceGUID
	}
	if gateway.ProvisioningState != "" {
		objectMap["provisioningState"] = gateway.ProvisioningState
	}
	if gateway.PublicIPAddresses != nil {
		objectMap["publicIpAddresses"] = gateway.PublicIPAddresses
	}
	if gateway.PublicIPPrefixes != nil {
		objectMap["publicIpPrefixes"] = gateway.PublicIPPrefixes
	}
	if gateway.Subnets != nil {
		objectMap["subnets"] = gateway.Subnets
	}

	if val, ok := objectMap[param]; ok {
		return val, nil
	}
	return nil, nil
}
