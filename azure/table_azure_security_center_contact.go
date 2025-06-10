package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
				Name:        "email",
				Description: "The email of this security contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Emails"),
			},
			{
				Name:        "phone",
				Description: "The phone number of this security contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Phone"),
			},
			{
				Name:        "is_enabled",
				Description: "Indicates whether the security contact is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.IsEnabled"),
			},
			{
				Name:        "alert_notifications",
				Description: "[DEPRECATED] Whether to send security alerts notifications to the security contact.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactProperties.AlertNotifications"),
			},
			{
				Name:        "alerts_to_admins",
				Description: "[DEPRECATED] Whether to send security alerts notifications to subscription admins.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContactProperties.AlertsToAdmins"),
			},
			{
				Name:        "notifications_by_role",
				Description: "Defines whether to send email notifications from Microsoft Defender for Cloud to persons with specific RBAC roles on the subscription.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.NotificationsByRole"),
			},
			{
				Name:        "notifications_sources",
				Description: "A collection of sources types which evaluate the email notification.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.NotificationsSources"),
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

func listSecurityCenterContacts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armsecurity.NewContactsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	pager := clientFactory.NewListPager(&armsecurity.ContactsClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, contact := range page.Value {
			d.StreamListItem(ctx, contact)

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

func getSecurityCenterContact(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armsecurity.NewContactsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQualString("name")

	result, err := clientFactory.Get(ctx, armsecurity.SecurityContactName(name), nil)
	if err != nil {
		return nil, err
	}

	return result.Contact, nil
}
