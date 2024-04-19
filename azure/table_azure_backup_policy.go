package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/backup"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/recoveryservices"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureBackupPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_backup_policy",
		Description: "Azure Backup Policy",
		List: &plugin.ListConfig{
			ParentHydrate: listRecoveryServicesVaults,
			Hydrate:       listBackupPolicy,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "vault_name",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_group",
					Require: plugin.Optional,
				},
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name associated with the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier representing the complete path to the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type representing the complete path of the form Namespace/ResourceType/ResourceType/...",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_name",
				Description: "The name of the vault associated with the backup policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags",
				Description: "The resource tags.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "azure_vm_workload_protection_policy_property",
				Description: "The Azure VM Workload Protection Policy associated with the backup policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureVMWorkloadProtectionPolicy"),
			},
			{
				Name:        "azure_file_share_protection_policy_property",
				Description: "The Azure File Share Protection Policy associated with the backup policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureFileShareProtectionPolicy"),
			},
			{
				Name:        "azure_iaas_vm_protection_policy_property",
				Description: "The Azure IaaS VM Protection Policy associated with the backup policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureIaaSVMProtectionPolicy"),
			},
			{
				Name:        "azure_sql_protection_policy_property",
				Description: "The Azure SQL Protection Policy associated with the backup policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AzureSqlProtectionPolicy"),
			},
			{
				Name:        "generic_protection_policy_property",
				Description: "The Generic Protection Policy associated with the backup policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("GenericProtectionPolicy"),
			},
			{
				Name:        "mab_protection_policy_property",
				Description: "The MAB Protection Policy associated with the backup policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MabProtectionPolicy"),
			},
			{
				Name:        "protection_policy_property",
				Description: "The protection policy associated with the backup policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProtectionPolicy"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
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

// // LIST FUNCTION
type ProtectionPolicyResource struct {
	backup.ProtectionPolicyResource
	VaultName                       string
	Location                        string
	AzureVMWorkloadProtectionPolicy *backup.AzureVMWorkloadProtectionPolicy
	AzureFileShareProtectionPolicy  *backup.AzureFileShareProtectionPolicy
	AzureIaaSVMProtectionPolicy     *backup.AzureIaaSVMProtectionPolicy
	AzureSqlProtectionPolicy        *backup.AzureSQLProtectionPolicy
	GenericProtectionPolicy         *backup.GenericProtectionPolicy
	MabProtectionPolicy             *backup.MabProtectionPolicy
	ProtectionPolicy                *backup.ProtectionPolicy
}

func listBackupPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vaultname := *h.Item.(recoveryservices.Vault).Name
	location := *h.Item.(recoveryservices.Vault).Location
	splitID := *h.Item.(recoveryservices.Vault).ID
	resourceGroupName := strings.Split(splitID, "/")[4]
	resourceGroupName = strings.ToLower(resourceGroupName)
	if d.EqualsQuals["vault_name"].GetStringValue() != "" || d.EqualsQuals["resource_group"].GetStringValue() != "" {
		if d.EqualsQuals["vault_name"].GetStringValue() != "" && d.EqualsQuals["vault_name"].GetStringValue() != vaultname {
			return nil, nil
		}
		if d.EqualsQuals["resource_group"].GetStringValue() != "" && d.EqualsQuals["resource_group"].GetStringValue() != resourceGroupName {
			return nil, nil
		}
	}
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	backupClient := backup.NewPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	backupClient.Authorizer = session.Authorizer

	result, err := backupClient.List(ctx, vaultname, resourceGroupName, "")
	if err != nil {
		return nil, err
	}
	for _, policy := range result.Values() {
		azureFileShareProtectionPolicy, _ := policy.Properties.AsAzureFileShareProtectionPolicy()
		azureVMWorkloadProtectionPolicy, _ := policy.Properties.AsAzureVMWorkloadProtectionPolicy()
		azureIaaSVMProtectionPolicy, _ := policy.Properties.AsAzureIaaSVMProtectionPolicy()
		azureSQLProtectionPolicy, _ := policy.Properties.AsAzureSQLProtectionPolicy()
		genericProtectionPolicy, _ := policy.Properties.AsGenericProtectionPolicy()
		mabProtectionPolicy, _ := policy.Properties.AsMabProtectionPolicy()
		protecttionPolicy, _ := policy.Properties.AsProtectionPolicy()
		backupPolicyOutput := ProtectionPolicyResource{
			policy,
			vaultname,
			location,
			azureVMWorkloadProtectionPolicy,
			azureFileShareProtectionPolicy,
			azureIaaSVMProtectionPolicy,
			azureSQLProtectionPolicy,
			genericProtectionPolicy,
			mabProtectionPolicy,
			protecttionPolicy,
		}
		d.StreamListItem(ctx, backupPolicyOutput)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, policy := range result.Values() {
			azureFileShareProtectionPolicy, _ := policy.Properties.AsAzureFileShareProtectionPolicy()
			azureVMWorkloadProtectionPolicy, _ := policy.Properties.AsAzureVMWorkloadProtectionPolicy()
			azureIaaSVMProtectionPolicy, _ := policy.Properties.AsAzureIaaSVMProtectionPolicy()
			azureSQLProtectionPolicy, _ := policy.Properties.AsAzureSQLProtectionPolicy()
			genericProtectionPolicy, _ := policy.Properties.AsGenericProtectionPolicy()
			mabProtectionPolicy, _ := policy.Properties.AsMabProtectionPolicy()
			protecttionPolicy, _ := policy.Properties.AsProtectionPolicy()
			backupPolicyOutput := ProtectionPolicyResource{
				policy,
				vaultname,
				location,
				azureVMWorkloadProtectionPolicy,
				azureFileShareProtectionPolicy,
				azureIaaSVMProtectionPolicy,
				azureSQLProtectionPolicy,
				genericProtectionPolicy,
				mabProtectionPolicy,
				protecttionPolicy,
			}
			d.StreamListItem(ctx, backupPolicyOutput)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}
