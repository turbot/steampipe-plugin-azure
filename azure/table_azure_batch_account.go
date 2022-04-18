package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2020-09-01/batch"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureBatchAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_batch_account",
		Description: "Azure Batch Account",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getBatchAccount,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "Invalid input"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listBatchAccounts,
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
				Name:        "provisioning_state",
				Description: "The provisioned state of the batch account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.ProvisioningState"),
			},
			{
				Name:        "account_endpoint",
				Description: "The account endpoint used to interact with the batch service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.AccountEndpoint"),
			},
			{
				Name:        "active_job_and_job_schedule_quota",
				Description: "Active job and job schedule quota of the batch account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.ActiveJobAndJobScheduleQuota"),
			},
			{
				Name:        "dedicated_core_quota",
				Description: "The dedicated core quota of the batch account.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AccountProperties.DedicatedCoreQuota"),
			},
			{
				Name:        "dedicated_core_quota_per_vm_family_enforced",
				Description: "Batch is transitioning its core quota system for dedicated cores to be enforced per Virtual Machine family. During this transitional phase, the dedicated core quota per Virtual Machine family may not yet be enforced.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AccountProperties.DedicatedCoreQuotaPerVMFamilyEnforced"),
			},
			{
				Name:        "low_priority_core_quota",
				Description: "The low priority core quota of the batch account.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AccountProperties.LowPriorityCoreQuota"),
			},
			{
				Name:        "pool_allocation_mode",
				Description: "The pool allocation mode of the batch account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.PoolAllocationMode"),
			},
			{
				Name:        "pool_quota",
				Description: "The pool quota of the batch account.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AccountProperties.PoolQuota"),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether or not public network access is allowed for the batch account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "auto_storage",
				Description: "The auto storage properties of the batch account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountProperties.AutoStorage"),
			},
			{
				Name:        "dedicated_core_quota_per_vm_family",
				Description: "A list of the dedicated core quota per virtual machine family for the batch account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountProperties.DedicatedCoreQuotaPerVMFamily"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the batch account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listBatchAccountDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "encryption",
				Description: "Properties to enable customer managed key for the batch account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountProperties.Encryption"),
			},
			{
				Name:        "identity",
				Description: "The identity of the batch account.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "key_vault_reference",
				Description: "Key vault reference of the batch account.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountProperties.KeyVaultReference"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "The properties associated with the private endpoint connection.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccountProperties.PrivateEndpointConnections"),
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

func listBatchAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	batchAccountClient := batch.NewAccountClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	batchAccountClient.Authorizer = session.Authorizer

	result, err := batchAccountClient.List(context.Background())
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

func getBatchAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBatchAccount")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	batchAccountClient := batch.NewAccountClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	batchAccountClient.Authorizer = session.Authorizer

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	op, err := batchAccountClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listBatchAccountDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listBatchAccountDiagnosticSettings")
	id := *h.Item.(batch.Account).ID

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
