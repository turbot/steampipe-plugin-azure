package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

//// TABLE DEFINITION

// The slots of the apps azure_app_service_function_app lists. The columns and the list hydrate are
// the web app slot table's: a slot is fetched the same way whatever kind of app it belongs to.
func tableAzureAppServiceFunctionAppSlot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_app_service_function_app_slot",
		Description: "Azure App Service Function App Slot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "app_name", "resource_group"}),
			Hydrate:    getAppServiceFunctionAppSlot,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "sites/slots/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAppServiceFunctionApps,
			Hydrate:       listAppServiceWebAppSlots,
			Tags: map[string]string{
				"service": "Microsoft.Web",
				"action":  "sites/slots/read",
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "app_name",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getConfigurationSlot,
				Tags: map[string]string{
					"service": "Microsoft.Web",
					"action":  "sites/slots/config/read",
				},
			},
		},
		Columns: azureColumns(appServiceSlotColumns()),
	}
}

//// HYDRATE FUNCTIONS

func getAppServiceFunctionAppSlot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	slot, err := getAppServiceSlot(ctx, d)
	// A slot has the same kind as its app, so only a function app's slots belong here
	if err != nil || slot == nil || !isFunctionAppKind(slot.Kind) {
		return nil, err
	}
	return slot, nil
}
