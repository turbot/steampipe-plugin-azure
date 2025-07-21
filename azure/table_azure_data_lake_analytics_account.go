package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/datalake/analytics/mgmt/account"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureDataLakeAnalyticsAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_data_lake_analytics_account",
		Description: "Azure Data Lake Analytics Account",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getDataLakeAnalyticsAccount,
			Tags: map[string]string{
				"service": "Microsoft.DataLakeAnalytics",
				"action":  "accounts/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listDataLakeAnalyticsAccounts,
			Tags: map[string]string{
				"service": "Microsoft.DataLakeAnalytics",
				"action":  "accounts/read",
			},
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
				Name:        "state",
				Description: "The state of the data lake analytics account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeAnalyticsAccountPropertiesBasic.State", "DataLakeAnalyticsAccountProperties.State"),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning status of the data lake analytics account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeAnalyticsAccountPropertiesBasic.ProvisioningState", "DataLakeAnalyticsAccountProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_id",
				Description: "The unique identifier associated with this data lake analytics account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeAnalyticsAccountPropertiesBasic.AccountID", "DataLakeAnalyticsAccountProperties.AccountID"),
			},
			{
				Name:        "creation_time",
				Description: "The data lake analytics account creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DataLakeAnalyticsAccountPropertiesBasic.CreationTime", "DataLakeAnalyticsAccountProperties.CreationTime").Transform(convertDateToTime),
			},
			{
				Name:        "current_tier",
				Description: "The commitment tier in use for current month.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.CurrentTier"),
			},
			{
				Name:        "default_data_lake_store_account",
				Description: "The default data lake store account associated with this data lake analytics account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.DefaultDataLakeStoreAccount"),
			},
			{
				Name:        "endpoint",
				Description: "The full cname endpoint for this data lake analytics account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLakeAnalyticsAccountPropertiesBasic.Endpoint", "DataLakeAnalyticsAccountProperties.Endpoint"),
			},
			{
				Name:        "firewall_state",
				Description: "The current state of the IP address firewall for this data lake analytics account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.FirewallState"),
			},
			{
				Name:        "last_modified_time",
				Description: "The data lake analytics account last modified time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DataLakeAnalyticsAccountPropertiesBasic.LastModifiedTime", "DataLakeAnalyticsAccountProperties.LastModifiedTime").Transform(convertDateToTime),
			},
			{
				Name:        "max_degree_of_parallelism",
				Description: "The maximum supported degree of parallelism for this data lake analytics account.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.MaxDegreeOfParallelism"),
			},
			{
				Name:        "max_degree_of_parallelism_per_job",
				Description: "The maximum supported degree of parallelism per job for this data lake analytics account.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.MaxDegreeOfParallelismPerJob"),
			},
			{
				Name:        "max_job_count",
				Description: "The maximum supported jobs running under the data lake analytics account at the same time.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.MaxJobCount"),
			},
			{
				Name:        "min_priority_per_job",
				Description: "The minimum supported priority per job for this data lake analytics account.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.MinPriorityPerJob"),
			},
			{
				Name:        "new_tier",
				Description: "The commitment tier to use for next month.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.NewTier"),
			},
			{
				Name:        "query_store_retention",
				Description: "The number of days that job metadata is retained.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.QueryStoreRetention"),
			},
			{
				Name:        "system_max_degree_of_parallelism",
				Description: "The system defined maximum supported degree of parallelism for this account, which restricts the maximum value of parallelism the user can set for the data lake analytics account.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.SystemMaxDegreeOfParallelism"),
			},
			{
				Name:        "system_max_job_count",
				Description: "The system defined maximum supported jobs running under the account at the same time, which restricts the maximum number of running jobs the user can set for the data lake analytics account.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.SystemMaxJobCount"),
			},
			{
				Name:        "compute_policies",
				Description: "The list of compute policies associated with this data lake analytics account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.ComputePolicies"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the data lake analytics account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listDataLakeAnalyticsAccountDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "data_lake_store_accounts",
				Description: "The list of data lake store accounts associated with this data lake analytics account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.DataLakeStoreAccounts"),
			},
			{
				Name:        "firewall_allow_azure_ips",
				Description: "The current state of allowing or disallowing IPs originating within azure through the firewall. If the firewall is disabled, this is not enforced.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.FirewallAllowAzureIps"),
			},
			{
				Name:        "firewall_rules",
				Description: "The list of firewall rules associated with this data lake analytics account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.FirewallRules"),
			},
			{
				Name:        "storage_accounts",
				Description: "The list of azure blob storage accounts associated with this data lake analytics account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDataLakeAnalyticsAccount,
				Transform:   transform.FromField("DataLakeAnalyticsAccountProperties.StorageAccounts"),
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

func listDataLakeAnalyticsAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := account.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &accountClient, d.Connection)

	result, err := accountClient.List(context.Background(), "", nil, nil, "", "", nil)
	if err != nil {
		return nil, err
	}
	for _, account := range result.Values() {
		d.StreamListItem(ctx, account)
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
			return nil, err
		}
		for _, account := range result.Values() {
			d.StreamListItem(ctx, account)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDataLakeAnalyticsAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDataLakeAnalyticsAccount")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	accountClient := account.NewAccountsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	accountClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &accountClient, d.Connection)

	var name, resourceGroup string
	if h.Item != nil {
		data := h.Item.(account.DataLakeAnalyticsAccountBasic)
		splitID := strings.Split(*data.ID, "/")
		name = *data.Name
		resourceGroup = splitID[4]
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		resourceGroup = d.EqualsQuals["resource_group"].GetStringValue()
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

func listDataLakeAnalyticsAccountDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDataLakeAnalyticsAccountDiagnosticSettings")
	id := getDataLakeAnalyticsAccountID(h.Item)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

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

func getDataLakeAnalyticsAccountID(item interface{}) string {
	switch item := item.(type) {
	case account.DataLakeAnalyticsAccountBasic:
		return *item.ID
	case account.DataLakeAnalyticsAccount:
		return *item.ID
	}
	return ""
}
