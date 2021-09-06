package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureProvider(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_provider",
		Description: "Azure Provider",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("namespace"),
			Hydrate:           getProvider,
			ShouldIgnoreError: isNotFoundError([]string{"InvalidResourceNamespace"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listProviders,
		},
		Columns: []*plugin.Column{
			{
				Name:        "namespace",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the resource provider.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a resource provider uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "registration_state",
				Description: "Contains the current registration state of the resource provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_types",
				Description: "A list of provider resource types.",
				Type:        proto.ColumnType_JSON,
			},

			// standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Namespace"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
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

//// LIST FUNCTION

func listProviders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	resourcesClient := resources.NewProvidersClient(subscriptionID)
	resourcesClient.Authorizer = session.Authorizer
	result, err := resourcesClient.List(ctx, nil, "")
	if err != nil {
		return nil, err
	}
	for _, provider := range result.Values() {
		d.StreamListItem(ctx, provider)
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if plugin.IsCancelled(ctx) {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, provider := range result.Values() {
			d.StreamListItem(ctx, provider)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getProvider")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	resourcesClient := resources.NewProvidersClient(subscriptionID)
	resourcesClient.Authorizer = session.Authorizer

	op, err := resourcesClient.Get(ctx, namespace, "")
	if err != nil {
		return nil, err
	}

	return op, nil
}
