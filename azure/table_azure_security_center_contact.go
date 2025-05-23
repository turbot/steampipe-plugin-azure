package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/security/mgmt/security"
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
		}),
	}
}

//// LIST FUNCTION

func listSecurityCenterContacts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	contactClient := security.NewContactsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	contactClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &contactClient, d.Connection)

	result, err := contactClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, contact := range result.Values() {
		d.StreamListItem(ctx, contact)
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
		for _, contact := range result.Values() {
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
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.EqualsQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	contactClient := security.NewContactsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	contactClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &contactClient, d.Connection)

	contact, err := contactClient.Get(ctx, name)
	if err != nil {
		return err, nil
	}

	return contact, nil
}
