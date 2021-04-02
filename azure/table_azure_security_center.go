package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2020-03-01-preview/policy"
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureSecurityCenter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center",
		Description: "Azure security Center",
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenter,
		},
		Columns: []*plugin.Column{
			{
				Name:        "setting",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pricing",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "auto_provisioning",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "contact",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("Security Center"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Id").Transform(idToAkas),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id").Transform(idToSubscriptionID),
			},
		},
	}
}

//// FETCH FUNCTIONS ////

func listSecurityCenter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	settingClient := security.NewSettingsClient(subscriptionID, "")
	settingClient.Authorizer = session.Authorizer

	pricingClient := security.NewPricingsClient(subscriptionID, "")
	pricingClient.Authorizer = session.Authorizer

	autoProvisionClient := security.NewAutoProvisioningSettingsClient(subscriptionID, "")
	autoProvisionClient.Authorizer = session.Authorizer

	contactClient := security.NewContactsClient(subscriptionID, "")
	contactClient.Authorizer = session.Authorizer

	PolicyClient := policy.NewAssignmentsClient()
	PolicyClient.Authorizer = session.Authorizer

	policy, err := PolicyClient.Get(ctx, "/subscriptions/"+subscriptionID, "SecurityCenterBuiltIn")
	if err != nil {
		return nil, err
	}
	contactList, err := contactClient.List(ctx)
	if err != nil {
		return nil, err
	}
	autoProvisionList, err := autoProvisionClient.List(ctx)
	if err != nil {
		return nil, err
	}
	pricingList, err := pricingClient.List(ctx)
	if err != nil {
		return nil, err
	}
	settingList, err := settingClient.List(ctx)
	if err != nil {
		return nil, err
	}
	var settings []map[string]interface{}
	var provisioning []map[string]interface{}
	var pricings []map[string]interface{}
	var contacts []map[string]interface{}

	for _, setting := range settingList.Values() {
		settingMap := make(map[string]interface{})
		settingMap["id"] = setting.ID
		settingMap["name"] = setting.Name
		settingMap["kind"] = setting.Kind
		settingMap["type"] = setting.Type
		settings = append(settings, settingMap)
	}

	for _, provision := range autoProvisionList.Values() {
		provisionMap := make(map[string]interface{})
		provisionMap["id"] = provision.ID
		provisionMap["name"] = provision.Name
		provisionMap["properties"] = provision.AutoProvisioningSettingProperties
		provisionMap["type"] = provision.Type
		provisioning = append(provisioning, provisionMap)
	}

	for _, pricing := range pricingList.Values() {
		pricingMap := make(map[string]interface{})
		pricingMap["id"] = pricing.ID
		pricingMap["name"] = pricing.Name
		pricingMap["properties"] = pricing.PricingProperties
		pricingMap["type"] = pricing.Type
		pricings = append(pricings, pricingMap)
	}

	for _, contact := range contactList.Values() {
		contactMap := make(map[string]interface{})
		contactMap["id"] = contact.ID
		contactMap["name"] = contact.Name
		contactMap["properties"] = contact.ContactProperties
		contactMap["type"] = contact.Type
		contacts = append(contacts, contactMap)
	}
	Id := "/subscriptions/" + subscriptionID + "/providers/Microsoft.Security/securityCenter"
	result := map[string]interface{}{
		"Setting":          settings,
		"Pricing":          pricings,
		"AutoProvisioning": provisioning,
		"Contact":          contacts,
		"Policy":           policy,
		"Id":               Id,
	}
	d.StreamListItem(ctx, result)
	return nil, err
}
