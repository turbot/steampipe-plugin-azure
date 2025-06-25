package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/policy"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzurePolicyDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_policy_definition",
		Description: "Azure Policy Definition",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPolicyDefinition,
			Tags: map[string]string{
				"service": "Microsoft.Authorization",
				"action":  "policyDefinitions/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listPolicyDefinitions,
			Tags: map[string]string{
				"service": "Microsoft.Authorization",
				"action":  "policyDefinitions/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
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
		}),
	}
}

//// LIST FUNCTION

func listPolicyDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	policyClient := policy.NewDefinitionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	policyClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &policyClient, d.Connection)

	result, err := policyClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, policy := range result.Values() {
		d.StreamListItem(ctx, policy)
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

		for _, policy := range result.Values() {
			d.StreamListItem(ctx, policy)
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

func getPolicyDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPolicyDefinition")

	name := d.EqualsQuals["name"].GetStringValue()

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_policy_definition.getPolicyDefinition", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := policy.NewDefinitionsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.Get(ctx, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_policy_definition.getPolicyDefinition", "api_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
