package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/policy"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzurePolicyAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_policy_assignment",
		Description: "Azure Policy Assignment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPolicyAssignment,
		},
		List: &plugin.ListConfig{
			Hydrate: listPolicyAssignments,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the policy assignment.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "The name of the policy assignment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the policy assignment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku_name",
				Description: "The name of the policy sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_tier",
				Description: "The policy sku tier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "assignment_properties",
				Description: "Properties for the policy assignment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "The managed identity associated with the policy assignment.",
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
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
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

func listPolicyAssignments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	PolicyClient := policy.NewAssignmentsClient(subscriptionID)
	PolicyClient.Authorizer = session.Authorizer

	policyList, err := PolicyClient.List(ctx, "")
	if err != nil {
		return err, nil
	}

	for _, policy := range policyList.Values() {
		d.StreamListItem(ctx, policy)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getPolicyAssignment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	PolicyClient := policy.NewAssignmentsClient(subscriptionID)
	PolicyClient.Authorizer = session.Authorizer

	policy, err := PolicyClient.Get(ctx, "/subscriptions/"+subscriptionID, name)
	if err != nil {
		return err, nil
	}

	return policy, nil
}
