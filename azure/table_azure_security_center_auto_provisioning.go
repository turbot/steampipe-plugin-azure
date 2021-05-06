package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterAutoProvisioning(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_auto_provisioning",
		Description: "Azure Security Center Auto Provisioning",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getSecurityCenterAutoProvisioning,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterAutoProvisioning,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The resource Id.",
				Transform:   transform.FromField("ID"),
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
				Description: "Describes what kind of security agent provisioning action to take. Possible values include: AutoProvisionOn, AutoProvisionOff",
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

func listSecurityCenterAutoProvisioning(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	autoProvisioningClient := security.NewAutoProvisioningSettingsClient(subscriptionID, "")
	autoProvisioningClient.Authorizer = session.Authorizer

	autoProvisioningList, err := autoProvisioningClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, autoProvisioning := range autoProvisioningList.Values() {
		d.StreamListItem(ctx, autoProvisioning)
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
	autoProvisioningClient := security.NewAutoProvisioningSettingsClient(subscriptionID, "")
	autoProvisioningClient.Authorizer = session.Authorizer

	autoProvisioning, err := autoProvisioningClient.Get(ctx, name)
	if err != nil {
		return err, nil
	}

	return autoProvisioning, nil
}
