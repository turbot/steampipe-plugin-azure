package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v3"
)

//// TABLE DEFINITION

func tableAzurePostgreSqlFlexibleServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_postgresql_flexible_server",
		Description: "Azure PostgreSQL Flexible Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getPostgreSqlFlexibleServer,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listPostgreSqlFlexibleServers,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "location",
				Description: "The geo-location where the resource lives.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity",
				Description: "Describes the identity of the application.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "properties",
				Description: "Properties of the server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sku",
				Description: "The SKU (pricing tier) of the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SKU"),
			},
			{
				Name:        "system_data",
				Description: "Azure Resource Manager metadata containing createdBy and modifiedBy information.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
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

func listPostgreSqlFlexibleServers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSqlFlexibleServers", "client_error", err)
		return nil, err
	}
	client := clientFactory.NewServersClient()
	resourceGroupName := h.Item.(resources.Group).Name

	result := client.NewListByResourceGroupPager(*resourceGroupName, &armpostgresqlflexibleservers.ServersClientListByResourceGroupOptions{})

	for result.More() {
	page, err := result.NextPage(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.listPostgreSqlFlexibleServers", "api_error", err)
		return nil, nil
	}

	for _, server := range page.Value {
		d.StreamListItem(ctx, server)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getPostgreSqlFlexibleServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPostgreSqlFlexibleServer")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.getPostgreSqlFlexibleServer", "credential_error", err)
	}

	name := d.EqualsQualString("name")
	resourceGroup := d.EqualsQualString("resource_group")

	// check if name or resourceGroup is empty
	if resourceGroup == "" || name == "" {
		return nil, nil
	}
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.getPostgreSqlFlexibleServer", "client_error", err)
	}

	res, err := clientFactory.NewServersClient().Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_postgresql_flexible_server.getPostgreSqlFlexibleServer", "api_error", err)
	}

	if res.Server.ID != nil {
		return res.Server, nil
	}

	return nil, nil
}
