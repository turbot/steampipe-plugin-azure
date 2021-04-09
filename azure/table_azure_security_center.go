package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/policy"
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center",
		Description: "Azure Security Center",
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenter,
		},
		Columns: []*plugin.Column{
			{
				Name:        "auto_provisioning",
				Description: "Auto provisioning settings of the subscriptions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "contact",
				Description: "Security contact configurations for the subscription.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy",
				Description: "Provides operations to assign policy definitions to a scope in your subscription.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pricing",
				Description: "Security pricing configuration in the resource group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "setting",
				Description: "Configuration settings for Azure Security Center.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
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

			// Azure standard columns
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listSecurityCenter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	// Fetch settings, provisioning, pricings, contacts and policy details from separate API's
	subscriptionID := session.SubscriptionID
	settings := getSettingDetails(ctx, session, subscriptionID)
	provisioning := getProvisioningDetails(ctx, session, subscriptionID)
	pricings := getPricingsDetails(ctx, session, subscriptionID)
	contacts := getContactDetails(ctx, session, subscriptionID)
	policy, err := getPolicyDetails(ctx, session, subscriptionID)

	if err != nil {
		return nil, err
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

func getSettingDetails(ctx context.Context, session *Session, subscriptionID string) []map[string]interface{} {
	settingClient := security.NewSettingsClient(subscriptionID, "")
	settingClient.Authorizer = session.Authorizer

	settingList, err := settingClient.List(ctx)
	if err != nil {
		return nil
	}
	var settings []map[string]interface{}

	for _, setting := range settingList.Values() {
		settingMap := make(map[string]interface{})
		settingMap["id"] = setting.ID
		settingMap["name"] = setting.Name
		settingMap["kind"] = setting.Kind
		settingMap["type"] = setting.Type
		settings = append(settings, settingMap)
	}
	return settings
}

func getProvisioningDetails(ctx context.Context, session *Session, subscriptionID string) []map[string]interface{} {
	autoProvisionClient := security.NewAutoProvisioningSettingsClient(subscriptionID, "")
	autoProvisionClient.Authorizer = session.Authorizer

	autoProvisionList, err := autoProvisionClient.List(ctx)
	if err != nil {
		return nil
	}

	var provisioning []map[string]interface{}

	for _, provision := range autoProvisionList.Values() {
		provisionMap := make(map[string]interface{})
		provisionMap["id"] = provision.ID
		provisionMap["name"] = provision.Name
		provisionMap["properties"] = provision.AutoProvisioningSettingProperties
		provisionMap["type"] = provision.Type
		provisioning = append(provisioning, provisionMap)
	}
	return provisioning
}

func getPricingsDetails(ctx context.Context, session *Session, subscriptionID string) []map[string]interface{} {
	pricingClient := security.NewPricingsClient(subscriptionID, "")
	pricingClient.Authorizer = session.Authorizer

	pricingList, err := pricingClient.List(ctx)
	if err != nil {
		return nil
	}

	var pricings []map[string]interface{}

	for _, pricing := range pricingList.Values() {
		pricingMap := make(map[string]interface{})
		pricingMap["id"] = pricing.ID
		pricingMap["name"] = pricing.Name
		pricingMap["properties"] = pricing.PricingProperties
		pricingMap["type"] = pricing.Type
		pricings = append(pricings, pricingMap)
	}
	return pricings
}

func getContactDetails(ctx context.Context, session *Session, subscriptionID string) []map[string]interface{} {
	contactClient := security.NewContactsClient(subscriptionID, "")
	contactClient.Authorizer = session.Authorizer

	contactList, err := contactClient.List(ctx)
	if err != nil {
		return nil
	}

	var contacts []map[string]interface{}

	for _, contact := range contactList.Values() {
		contactMap := make(map[string]interface{})
		contactMap["id"] = contact.ID
		contactMap["name"] = contact.Name
		contactMap["properties"] = contact.ContactProperties
		contactMap["type"] = contact.Type
		contacts = append(contacts, contactMap)
	}
	return contacts
}

func getPolicyDetails(ctx context.Context, session *Session, subscriptionID string) (policy.Assignment, error) {
	PolicyClient := policy.NewAssignmentsClient(subscriptionID)
	PolicyClient.Authorizer = session.Authorizer

	policy, err := PolicyClient.Get(ctx, "/subscriptions/"+subscriptionID, "SecurityCenterBuiltIn")

	return policy, err
}
