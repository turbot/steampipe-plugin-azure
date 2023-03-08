package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterAutomation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_automation",
		Description: "Azure Security Center Automation",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getSecurityCenterAutomation,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterAutomations,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The resource id.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The security automation description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AutomationProperties.Description"),
			},
			{
				Name:        "is_enabled",
				Description: "Indicates whether the security automation is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AutomationProperties.IsEnabled"),
			},
			{
				Name:        "kind",
				Description: "Kind of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "Entity tag is used for comparing two or more entities from the same requested resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actions",
				Description: "A collection of the actions which are triggered if all the configured rules evaluations, within at least one rule set, are true.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AutomationProperties.Actions"),
			},
			{
				Name:        "scopes",
				Description: "A collection of scopes on which the security automations logic is applied. Supported scopes are the subscription itself or a resource group under that subscription. The automation will only apply on defined scopes.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AutomationProperties.Scopes"),
			},
			{
				Name:        "sources",
				Description: "A collection of the source event types which evaluate the security automation set of rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AutomationProperties.Sources"),
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
				Description: "A list of key value pairs that describe the resource.",
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

func listSecurityCenterAutomations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	automationClient := security.NewAutomationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	automationClient.Authorizer = session.Authorizer

	result, err := automationClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, automation := range result.Values() {
		d.StreamListItem(ctx, automation)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return err, nil
		}

		for _, automation := range result.Values() {
			d.StreamListItem(ctx, automation)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCenterAutomation(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	automationClient := security.NewAutomationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	automationClient.Authorizer = session.Authorizer

	automation, err := automationClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return err, nil
	}

	return automation, nil
}
