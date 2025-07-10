package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/security/mgmt/security"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_setting",
		Description: "Azure Security Center Setting",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getSecurityCenterSetting,
			Tags: map[string]string{
				"service": "Microsoft.Security",
				"action":  "settings/read",
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterSettings,
			Tags: map[string]string{
				"service": "Microsoft.Security",
				"action":  "settings/read",
			},
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
				Name:        "enabled",
				Description: "Check if the setting is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The kind of the setting.",
				Type:        proto.ColumnType_STRING,
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

type SecurityCenterSettings struct {
	ID      *string
	Name    *string
	Enabled *bool
	Type    *string
	Kind    security.KindEnum2
}

//// LIST FUNCTION

func listSecurityCenterSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	settingClient := security.NewSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	settingClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &settingClient, d.Connection)

	result, err := settingClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, setting := range result.Values() {

		// Check if dataExportSettings
		if dataExportSettings, ok := setting.AsDataExportSettings(); ok {
			d.StreamListItem(ctx, SecurityCenterSettings{
				ID:      dataExportSettings.ID,
				Name:    dataExportSettings.Name,
				Enabled: dataExportSettings.Enabled,
				Type:    dataExportSettings.Type,
				Kind:    dataExportSettings.Kind,
			})
		} else if alertSyncSettings, ok := setting.AsAlertSyncSettings(); ok { // Check if alertSyncSettings
			d.StreamListItem(ctx, SecurityCenterSettings{
				ID:      alertSyncSettings.ID,
				Name:    alertSyncSettings.Name,
				Enabled: alertSyncSettings.Enabled,
				Type:    alertSyncSettings.Type,
				Kind:    alertSyncSettings.Kind,
			})
		} else {
			setting, _ := setting.AsSetting() // Basic settings
			d.StreamListItem(ctx, SecurityCenterSettings{
				ID:   setting.ID,
				Name: setting.Name,
				Type: setting.Type,
				Kind: setting.Kind,
			})
		}

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
			return err, nil
		}
		for _, setting := range result.Values() {

			// Check if dataExportSettings
			if dataExportSettings, ok := setting.AsDataExportSettings(); ok {
				d.StreamListItem(ctx, SecurityCenterSettings{
					ID:      dataExportSettings.ID,
					Name:    dataExportSettings.Name,
					Enabled: dataExportSettings.Enabled,
					Type:    dataExportSettings.Type,
					Kind:    dataExportSettings.Kind,
				})
			} else if alertSyncSettings, ok := setting.AsAlertSyncSettings(); ok { // Check if alertSyncSettings
				d.StreamListItem(ctx, SecurityCenterSettings{
					ID:      alertSyncSettings.ID,
					Name:    alertSyncSettings.Name,
					Enabled: alertSyncSettings.Enabled,
					Type:    alertSyncSettings.Type,
					Kind:    alertSyncSettings.Kind,
				})
			} else {
				setting, _ := setting.AsSetting() // Basic settings
				d.StreamListItem(ctx, SecurityCenterSettings{
					ID:   setting.ID,
					Name: setting.Name,
					Type: setting.Type,
					Kind: setting.Kind,
				})
			}

			// Check if context has been cancelled or if the limit has been hit (if specified)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCenterSetting(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.EqualsQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	settingClient := security.NewSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	settingClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &settingClient, d.Connection)

	setting, err := settingClient.Get(ctx, security.SettingName4(name))
	if err != nil {
		return err, nil
	}

	if dataExportSettings, ok := setting.Value.AsDataExportSettings(); ok {
		return dataExportSettings, nil
	} else if alertSyncSettings, ok := setting.Value.AsAlertSyncSettings(); ok {
		return alertSyncSettings, nil
	} else {
		settings, _ := setting.Value.AsSetting()
		return settings, nil
	}
}
