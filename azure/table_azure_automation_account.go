package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/automation/mgmt/automation"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureApAutomationAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_automation_account",
		Description: "Azure Automation Account",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAutomationAccount,
			Tags: map[string]string{
				"service": "Microsoft.Automation",
				"action":  "automationAccounts/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAutomationAccounts,
			Tags: map[string]string{
				"service": "Microsoft.Automation",
				"action":  "automationAccounts/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource.",
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "description",
				Description: "The description for the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.Description"),
			},
			{
				Name:        "etag",
				Description: "Gets the etag of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the account.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AccountProperties.CreationTime.Time"),
			},
			{
				Name:        "last_modified_time",
				Description: "The last modified time of the account.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AccountProperties.LastModifiedTime.Time"),
			},
			{
				Name:        "last_modified_by",
				Description: "The account last modified by.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.LastModifiedBy"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The status of account. Possible values include: 'AccountStateOk', 'AccountStateUnavailable', 'AccountStateSuspended'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.State"),
			},
			{
				Name:        "sku_name",
				Description: "The SKU name of the account. Possible values include: 'SkuNameEnumFree', 'SkuNameEnumBasic'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_family",
				Description: "The SKU family of the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Family"),
			},
			{
				Name:        "sku_capacity",
				Description: "The SKU capacity of the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Capacity"),
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

//// LIST FUNCTION ////

func listAutomationAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.listAutomationAccounts", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := automation.NewAccountClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &accountClient, d.Connection)

	result, err := accountClient.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.listAutomationAccounts", "api_error", err)
		return nil, err
	}

	for _, account := range result.Values() {
		d.StreamListItem(ctx, account)
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
			plugin.Logger(ctx).Error("azure_automation_variable.listAutomationAccounts", "paginator_error", err)
			return nil, err
		}

		for _, account := range result.Values() {
			d.StreamListItem(ctx, account)
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

func getAutomationAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.getAutomationAccount", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := automation.NewAccountClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &accountClient, d.Connection)

	op, err := accountClient.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.getAutomationAccount", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// Instead it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
