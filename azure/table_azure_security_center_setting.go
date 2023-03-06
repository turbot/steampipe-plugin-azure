package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
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
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterSettings,
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
				Description: "Data export setting status.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DataExportSettingProperties.Enabled"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The kind of the settings string (DataExportSettings).",
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

//// LIST FUNCTION

func listSecurityCenterSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	settingClient := security.NewSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	settingClient.Authorizer = session.Authorizer

	result, err := settingClient.List(ctx)
	if err != nil {
		return err, nil
	}

	for _, setting := range result.Values() {
		d.StreamListItem(ctx, setting)
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
		for _, setting := range result.Values() {
			d.StreamListItem(ctx, setting)
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

func getSecurityCenterSetting(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	name := d.EqualsQuals["name"].GetStringValue()

	subscriptionID := session.SubscriptionID
	settingClient := security.NewSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	settingClient.Authorizer = session.Authorizer

	setting, err := settingClient.Get(ctx, name)
	if err != nil {
		return err, nil
	}

	return setting.Value, nil
}
