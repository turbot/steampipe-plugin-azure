package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/security/mgmt/security"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterJITNetworkAccessPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_jit_network_access_policy",
		Description: "Azure Security Center JIT Network Access Policy",
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterJITNetworkAccessPolicies,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Kind of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the Just-in-Time policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JitNetworkAccessPolicyProperties.ProvisioningState"),
			},
			{
				Name:        "virtual_machines",
				Description: "Configurations for Microsoft.Compute/virtualMachines resource type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("JitNetworkAccessPolicyProperties.VirtualMachines"),
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
		}),
	}
}

//// LIST FUNCTION

func listSecurityCenterJITNetworkAccessPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	client := security.NewJitNetworkAccessPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, jitNetworkAccessPolicy := range result.Values() {
		d.StreamListItem(ctx, jitNetworkAccessPolicy)
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
		for _, jitNetworkAccessPolicy := range result.Values() {
			d.StreamListItem(ctx, jitNetworkAccessPolicy)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
