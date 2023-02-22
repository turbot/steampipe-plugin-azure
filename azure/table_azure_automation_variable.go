package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2019-06-01/automation"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION ////

func tableAzureApAutomationVariable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_automation_variable",
		Description: "Azure Automation Variable",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"account_name", "name", "resource_group"}),
			Hydrate:    getAutomationVariable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAutomationAccounts,
			Hydrate:       listAutomationVariables,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource.",
			},
			{
				Name:        "account_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the account.",
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "description",
				Description: "The description for the variable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VariableProperties.Description"),
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the variable.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("VariableProperties.CreationTime.Time"),
			},
			{
				Name:        "last_modified_time",
				Description: "The last modified time of the variable.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("VariableProperties.LastModifiedTime.Time"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_encrypted",
				Description: "The encrypted flag of the variable.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("VariableProperties.IsEncrypted"),
			},
			{
				Name:        "value",
				Description: "The value of the variable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VariableProperties.Value"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

type VariableDetails struct {
	AccountName string
	automation.Variable
}

//// LIST FUNCTION ////

func listAutomationVariables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.listAutomationVariables", "session_error", err)
		return nil, err
	}

	var account automation.Account
	if h.Item != nil {
		account = h.Item.(automation.Account)
	} else {
		return nil, nil
	}
	resourceGroupName := strings.Split(*account.ID, "/")[4]
	accountName := account.Name

	subscriptionID := session.SubscriptionID

	accountClient := automation.NewVariableClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	result, err := accountClient.ListByAutomationAccount(ctx, resourceGroupName, *accountName)
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.listAutomationVariables", "api_error", err)
		return nil, err
	}

	for _, variable := range result.Values() {
		d.StreamListItem(ctx, &VariableDetails{*accountName, variable})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_automation_variable.listAutomationVariables", "paginator_error", err)
			return nil, err
		}

		for _, variable := range result.Values() {
			d.StreamListItem(ctx, &VariableDetails{*accountName, variable})
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getAutomationVariable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	accountName := d.KeyColumnQuals["account_name"].GetStringValue()
	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.getAutomationVariable", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := automation.NewVariableClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	op, err := accountClient.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_automation_variable.getAutomationVariable", "api_error", err)
		return nil, err
	}

	// In some cases the API does not return any notFound error
	// instead it returns empty data
	if op.ID != nil {
		return &VariableDetails{accountName, op}, nil
	}

	return nil, nil
}
