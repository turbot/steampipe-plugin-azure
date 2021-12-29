package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureDataLakeStore(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_lake_store",
		Description: "Azure Data Lake Store",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getDataLakeStore,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listDataLakeStores,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_id",
				Description: "The unique identifier associated with this data lake store account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.AccountID", "DataLakeStoreAccountProperties.AccountID"),
			},
			{
				Name:        "creation_time",
				Description: "The account creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.CreationTime", "DataLakeStoreAccountProperties.CreationTime").Transform(convertDateToTime),
			},
			{
				Name:        "current_tier",
				Description: "The commitment tier in use for current month.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.CurrentTier"),
			},
			{
				Name:        "default_group",
				Description: "The default owner group for all new folders and files created in the data lake store account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.DefaultGroup"),
			},
			{
				Name:        "encryption_state",
				Description: "The current state of encryption for this data lake store account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.EncryptionState"),
			},
			{
				Name:        "encryption_provisioning_state",
				Description: "The current state of encryption provisioning for this data lake store account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.EncryptionProvisioningState"),
			},
			{
				Name:        "endpoint",
				Description: "The full cname endpoint for this account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.Endpoint", "DataLakeStoreAccountProperties.Endpoint"),
			},
			{
				Name:        "firewall_state",
				Description: "The current state of the IP address firewall for this data lake store account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.FirewallState"),
			},
			{
				Name:        "last_modified_time",
				Description: "The account last modified time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.LastModifiedTime", "DataLakeStoreAccountProperties.LastModifiedTime").Transform(convertDateToTime),
			},
			{
				Name:        "new_tier",
				Description: "The commitment tier to use for next month.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.NewTier"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning status of the data lake store account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.ProvisioningState", "DataLakeStoreAccountProperties.ProvisioningState"),
			},
			{
				Name:        "state",
				Description: "The state of the data lake store account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeStoreAccountPropertiesBasic.State", "DataLakeStoreAccountProperties.State"),
			},
			{
				Name:        "trusted_id_provider_state",
				Description: "The current state of the trusted identity provider feature for this data lake store account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.TrustedIDProviderState"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the data lake store.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listDataLakeStoreDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption_config",
				Description: "The key vault encryption configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.EncryptionConfig"),
			},
			{
				Name:        "firewall_allow_azure_ips",
				Description: "The current state of allowing or disallowing IPs originating within azure through the firewall. If the firewall is disabled, this is not enforced.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.FirewallAllowAzureIps"),
			},
			{
				Name:        "firewall_rules",
				Description: "The list of firewall rules associated with this data lake store account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.FirewallRules"),
			},
			{
				Name:        "identity",
				Description: "The key vault encryption identity, if any.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeStore,
			},
			{
				Name:        "trusted_id_providers",
				Description: "The list of trusted identity providers associated with this data lake store account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.TrustedIDProviders"),
			},
			{
				Name:        "virtual_network_rules",
				Description: "The list of virtual network rules associated with this data lake store account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeStore,
				Transform:   transform.FromField("DataLakeStoreAccountProperties.VirtualNetworkRules"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listDataLakeStores(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := account.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	result, err := accountClient.List(ctx, "", nil, nil, "", "", nil)
	if err != nil {
		return nil, err
	}
	for _, account := range result.Values() {
		d.StreamListItem(ctx, account)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, account := range result.Values() {
			d.StreamListItem(ctx, account)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDataLakeStore(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataLakeStore")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := account.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	var name, resourceGroup string
	if h.Item != nil {
		data := h.Item.(account.DataLakeStoreAccountBasic)
		splitID := strings.Split(*data.ID, "/")
		name = *data.Name
		resourceGroup = splitID[4]
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		resourceGroup = d.KeyColumnQuals["resource_group"].GetStringValue()
	}

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := accountClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listDataLakeStoreDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDataLakeStoreDiagnosticSettings")
	id := getLakeStoreID(h.Item)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives
	// the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["id"] = i.ID
		}
		if i.Name != nil {
			objectMap["name"] = i.Name
		}
		if i.Type != nil {
			objectMap["type"] = i.Type
		}
		if i.DiagnosticSettings != nil {
			objectMap["properties"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}
	return diagnosticSettings, nil
}

func getLakeStoreID(item interface{}) string {
	switch item := item.(type) {
	case account.DataLakeStoreAccountBasic:
		return *item.ID
	case account.DataLakeStoreAccount:
		return *item.ID
	}
	return ""
}
