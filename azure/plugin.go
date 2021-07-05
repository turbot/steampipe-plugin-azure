package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-azure"

// Plugin creates this (azure) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound"}),
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"azure_ad_group":                             tableAzureAdGroup(ctx),
			"azure_ad_service_principal":                 tableAzureAdServicePrincipal(ctx),
			"azure_ad_user":                              tableAzureAdUser(ctx),
			"azure_api_management":                       tableAzureAPIManagement(ctx),
			"azure_app_service_environment":              tableAzureAppServiceEnvironment(ctx),
			"azure_app_service_function_app":             tableAzureAppServiceFunctionApp(ctx),
			"azure_app_service_plan":                     tableAzureAppServicePlan(ctx),
			"azure_app_service_web_app":                  tableAzureAppServiceWebApp(ctx),
			"azure_application_security_group":           tableAzureApplicationSecurityGroup(ctx),
			"azure_compute_availability_set":             tableAzureComputeAvailabilitySet(ctx),
			"azure_compute_disk":                         tableAzureComputeDisk(ctx),
			"azure_compute_disk_encryption_set":          tableAzureComputeDiskEncryptionSet(ctx),
			"azure_compute_image":                        tableAzureComputeImage(ctx),
			"azure_compute_resource_sku":                 tableAzureResourceSku(ctx),
			"azure_compute_snapshot":                     tableAzureComputeSnapshot(ctx),
			"azure_compute_virtual_machine":              tableAzureComputeVirtualMachine(ctx),
			"azure_cosmosdb_account":                     tableAzureCosmosDBAccount(ctx),
			"azure_cosmosdb_mongo_database":              tableAzureCosmosDBMongoDatabase(ctx),
			"azure_cosmosdb_sql_database":                tableAzureCosmosDBSQLDatabase(ctx),
			"azure_data_factory":                         tableAzureDataFactory(ctx),
			"azure_data_factory_pipeline":                tableAzureDataFactoryPipeline(ctx),
			"azure_diagnostic_setting":                   tableAzureDiagnosticSetting(ctx),
			"azure_firewall":                             tableAzureFirewall(ctx),
			"azure_key_vault":                            tableAzureKeyVault(ctx),
			"azure_key_vault_key":                        tableAzureKeyVaultKey(ctx),
			"azure_key_vault_secret":                     tableAzureKeyVaultSecret(ctx),
			"azure_kubernetes_cluster":                   tableAzureKubernetesCluster(ctx),
			"azure_location":                             tableAzureLocation(ctx),
			"azure_log_alert":                            tableAzureLogAlert(ctx),
			"azure_log_profile":                          tableAzureLogProfile(ctx),
			"azure_management_lock":                      tableAzureManagementLock(ctx),
			"azure_mysql_server":                         tableAzureMySQLServer(ctx),
			"azure_network_interface":                    tableAzureNetworkInterface(ctx),
			"azure_network_security_group":               tableAzureNetworkSecurityGroup(ctx),
			"azure_network_watcher":                      tableAzureNetworkWatcher(ctx),
			"azure_network_watcher_flow_log":             tableAzureNetworkWatcherFlowLog(ctx),
			"azure_policy_assignment":                    tableAzurePolicyAssignment(ctx),
			"azure_policy_definition":                    tableAzurePolicyDefinition(ctx),
			"azure_postgresql_server":                    tableAzurePostgreSqlServer(ctx),
			"azure_provider":                             tableAzureProvider(ctx),
			"azure_public_ip":                            tableAzurePublicIP(ctx),
			"azure_resource_group":                       tableAzureResourceGroup(ctx),
			"azure_role_assignment":                      tableAzureIamRoleAssignment(ctx),
			"azure_role_definition":                      tableAzureIamRoleDefinition(ctx),
			"azure_route_table":                          tableAzureRouteTable(ctx),
			"azure_security_center_auto_provisioning":    tableAzureSecurityCenterAutoProvisioning(ctx),
			"azure_security_center_contact":              tableAzureSecurityCenterContact(ctx),
			"azure_security_center_setting":              tableAzureSecurityCenterSetting(ctx),
			"azure_security_center_subscription_pricing": tableAzureSecurityCenterPricing(ctx),
			"azure_sql_database":                         tableAzureSqlDatabase(ctx),
			"azure_sql_server":                           tableAzureSQLServer(ctx),
			"azure_storage_account":                      tableAzureStorageAccount(ctx),
			"azure_storage_blob":                         tableAzureStorageBlob(ctx),
			"azure_storage_blob_service":                 tableAzureStorageBlobService(ctx),
			"azure_storage_container":                    tableAzureStorageContainer(ctx),
			"azure_storage_queue":                        tableAzureStorageQueue(ctx),
			"azure_storage_table":                        tableAzureStorageTable(ctx),
			"azure_storage_table_service":                tableAzureStorageTableService(ctx),
			"azure_subnet":                               tableAzureSubnet(ctx),
			"azure_subscription":                         tableAzureSubscription(ctx),
			"azure_tenant":                               tableAzureTenant(ctx),
			"azure_virtual_network":                      tableAzureVirtualNetwork(ctx),
			// "azure_storage_table":               tableAzureStorageTable(ctx),
		},
	}

	return p
}
