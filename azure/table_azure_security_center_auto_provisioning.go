package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterAutoProvisioning(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_auto_provisioning",
		Description: "Azure Security Center Auto Provisioning",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getSecurityCenterAutoProvisioning,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterAutoProvisioning,
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
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_provision",
				Description: "Describes what kind of security agent provisioning action to take. Possible values include: On, Off",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AutoProvisioningSettingProperties.AutoProvision"),
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

func listSecurityCenterAutoProvisioning(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	autoProvisioningClient := security.NewAutoProvisioningSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	autoProvisioningClient.Authorizer = session.Authorizer

	result, err := autoProvisioningClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, autoProvisioning := range result.Values() {
		d.StreamListItem(ctx, autoProvisioning)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return err, nil
		}
		for _, contact := range result.Values() {
			d.StreamListItem(ctx, contact)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCenterAutoProvisioning(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	autoProvisioningClient := security.NewAutoProvisioningSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	autoProvisioningClient.Authorizer = session.Authorizer

	autoProvisioning, err := autoProvisioningClient.Get(ctx, name)
	if err != nil {
		return err, nil
	}

	return autoProvisioning, nil
}
