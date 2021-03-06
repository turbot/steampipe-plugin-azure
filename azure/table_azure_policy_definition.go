package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzurePolicyDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_policy_definition",
		Description: "Azure Policy Definition",
		// Get API operation is not working as expected, skipping for now
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.SingleColumn("name"),
		// 	Hydrate:    getPolicyDefinition,
		// },
		List: &plugin.ListConfig{
			Hydrate: listPolicyDefintions,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the policy definition.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "The name of the policy definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "display_name",
				Description: "The user-friendly display name of the policy definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefinitionProperties.DisplayName"),
			},
			{
				Name:        "description",
				Description: "The policy definition description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefinitionProperties.Description"),
			},
			{
				Name:        "mode",
				Description: "The policy definition mode. Some examples are All, Indexed, Microsoft.KeyVault.Data.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefinitionProperties.Mode"),
			},
			{
				Name:        "policy_type",
				Description: "The type of policy definition. Possible values are NotSpecified, BuiltIn, Custom, and Static. Possible values include: 'NotSpecified', 'BuiltIn', 'Custom', 'Static'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefinitionProperties.PolicyType"),
			},
			{
				Name:        "type",
				Description: "The type of the resource (Microsoft.Authorization/policyDefinitions).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "metadata",
				Description: "The policy definition metadata.  Metadata is an open ended object and is typically a collection of key value pairs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefinitionProperties.Metadata"),
			},
			{
				Name:        "parameters",
				Description: "The parameter definitions for parameters used in the policy rule. The keys are the parameter names.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefinitionProperties.Parameters"),
			},
			{
				Name:        "policy_rule",
				Description: "The policy rule.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefinitionProperties.PolicyRule"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefinitionProperties.DisplayName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPolicyDefinitionTurbotData,
			},

			// Azure standard columns
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getPolicyDefinitionTurbotData,
			},
		},
	}
}

//// LIST FUNCTION

func listPolicyDefintions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	PolicyClient := policy.NewDefinitionsClient(subscriptionID)
	PolicyClient.Authorizer = session.Authorizer

	policyList, err := PolicyClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, policy := range policyList.Values() {
		d.StreamListItem(ctx, policy)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getPolicyDefinitionTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPolicyDefinitionTurbotData")
	data := h.Item.(policy.Definition)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	akas := []string{"azure:///subscriptions/" + subscriptionID + *data.ID, "azure:///subscriptions/" + subscriptionID + strings.ToLower(*data.ID)}

	turbotData := map[string]interface{}{
		"SubscriptionId": subscriptionID,
		"Akas":           akas,
	}

	return turbotData, nil
}
