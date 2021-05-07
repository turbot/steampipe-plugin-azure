package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterContact(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_contact",
		Description: "Azure Security Center Contact",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getSecurityCenterContact,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterContacts,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The resource id.",
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
				Name:        "email",
				Description: "The email of this security contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactProperties.Email"),
			},
			{
				Name:        "phone",
				Description: "The phone number of this security contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactProperties.Phone"),
			},
			{
				Name:        "alert_notifications",
				Description: "Whether to send security alerts notifications to the security contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactProperties.AlertNotifications"),
			},
			{
				Name:        "alerts_to_admins",
				Description: "Whether to send security alerts notifications to subscription admins.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactProperties.AlertsToAdmins"),
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

func listSecurityCenterContacts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	contactClient := security.NewContactsClient(subscriptionID, "")
	contactClient.Authorizer = session.Authorizer

	contactList, err := contactClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, contact := range contactList.Values() {
		d.StreamListItem(ctx, contact)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCenterContact(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	contactClient := security.NewContactsClient(subscriptionID, "")
	contactClient.Authorizer = session.Authorizer

	contact, err := contactClient.Get(ctx, name)
	if err != nil {
		return err, nil
	}

	return contact, nil
}
