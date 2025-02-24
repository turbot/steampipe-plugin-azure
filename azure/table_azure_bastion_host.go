package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureBastionHost(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_bastion_host",
		Description: "Azure Bastion Host",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getBastionHost,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listBastionHosts,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the bastion host.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a bastion host uniquely.",
				Transform:   transform.FromGo(),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_name",
				Description: "FQDN for the endpoint on which bastion host is accessible.",
				Transform:   transform.FromField("BastionHostPropertiesFormat.DNSName"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "A unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the bastion host resource.",
				Transform:   transform.FromField("BastionHostPropertiesFormat.ProvisioningState"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the bastion host.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_configurations",
				Description: "IP configuration of the bastion host resource.",
				Transform:   transform.FromField("BastionHostPropertiesFormat.IPConfigurations"),
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Transform:   transform.FromField("ID").Transform(idToAkas),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Name"),
				Type:        proto.ColumnType_STRING,
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Transform:   transform.FromField("Location").Transform(toLower),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// FETCH FUNCTIONS ////

func listBastionHosts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_bastion_host.listBastionHosts", "client_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	bastionClient := network.NewBastionHostsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	bastionClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &bastionClient, d.Connection)

	result, err := bastionClient.List(ctx)
	if err != nil {
		logger.Error("azure_bastion_host.listBastionHosts", "api_error", err)
		return nil, err
	}

	for _, host := range result.Values() {
		d.StreamListItem(ctx, host)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			logger.Error("azure_bastion_host.listBastionHosts", "api_error", err)
			return nil, err
		}

		for _, host := range result.Values() {
			d.StreamListItem(ctx, host)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getBastionHost(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_bastion_host.getBastionHost", "client_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	bastionClient := network.NewBastionHostsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	bastionClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &bastionClient, d.Connection)

	result, err := bastionClient.Get(ctx, resourceGroup, name)
	if err != nil {
		logger.Error("azure_bastion_host.getBastionHost", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if result.ID != nil {
		return result, nil
	}

	return nil, nil
}
