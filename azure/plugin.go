package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/rate_limiter"
)

const pluginName = "steampipe-plugin-azure"

// Plugin creates this (azure) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		DefaultGetConfig: &plugin.GetConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound"}),
			},
		},
		// Default ignore config for the plugin
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrorPluginDefault(),
		},
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "subscription_id",
				Hydrate: getSubscriptionIdForConnection,
			},
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		RateLimiters: []*rate_limiter.Definition{
			// Tables mentioned in GitHub issue #927 - Azure Services Highly Prone to Rate Limits

			// 1. Azure Resource Manager (ARM) - All list and get resource metadata calls go through ARM
			// https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/request-limits-and-throttling
			// ~12,000 reads/hour per subscription per region
			{
				Name:       "azure_subscription",
				FillRate:   25,
				BucketSize: 250,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'Microsoft.Resources' and action = 'subscriptions/read'",
			},
			// https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/request-limits-and-throttling#subscription-and-tenant-limits
			{
				Name:       "azure_resource_group",
				FillRate:   25,
				BucketSize: 250,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'Microsoft.Resources' and action = 'resourceGroups/read'",
			},
			// https://learn.microsoft.com/en-us/azure/virtual-machines/compute-throttling-limits#throttling-limits-for-virtual-machines
			{
				Name:       "azure_compute_virtual_machine",
				FillRate:   500,
				BucketSize: 1500,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'Microsoft.Compute' and action = 'virtualMachines/read'",
			},
			{
				Name:       "azure_compute_virtual_machine_operations",
				FillRate:   12,
				BucketSize: 36,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service in ('Microsoft.Compute', 'Microsoft.GuestConfiguration', 'Microsoft.Network') and action in ('virtualMachines/instanceView/read', 'virtualMachines/extensions/read', 'publicIPAddresses/read')",
			},

			// https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/azure-subscription-service-limits#azure-storage-resource-provider-limits
			{
				Name:       "azure_storage_account",
				FillRate:   2,
				BucketSize: 50,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'Microsoft.Storage' and action = 'storageAccounts/read'",
			},
			// https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/azure-subscription-service-limits#azure-blob-storage-limits
			{
				Name:       "azure_storage_blob",
				FillRate:   500,
				BucketSize: 500,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'Microsoft.Storage' and action = 'storageAccounts/blobServices/containers/blobs/read'",
			},

			// 3. Azure Key Vault - Every get for a secret/key/certificate is a counted API request
			// https://learn.microsoft.com/en-us/azure/key-vault/general/service-limits
			// ~2,000 GET/10s per vault
			{
				Name:       "azure_key_vault_secret",
				FillRate:   40,
				BucketSize: 400,
				Scope:      []string{"connection", "subscription", "vault"},
				Where:      "service = 'Microsoft.KeyVault' and action = 'vaults/secrets/read'",
			},
			{
				Name:       "azure_key_vault_key",
				FillRate:   20,
				BucketSize: 200,
				Scope:      []string{"connection", "subscription", "vault"},
				Where:      "service = 'Microsoft.KeyVault' and action = 'vaults/keys/read'",
			},

			// Azure Monitor / Log Analytics - Querying log data is essentially a read op
			// https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/request-limits-and-throttling
			// API queries often capped by 200 QPS per workspace, but contributes significantly to ARM traffic
			// https://learn.microsoft.com/en-us/azure/azure-monitor/fundamentals/service-limits#alerts-api
			{
				Name:       "azure_monitor_activity_log",
				FillRate:   1, // Conservative limit for activity log queries
				BucketSize: 50,
				Scope:      []string{"connection", "subscription", "region"},
				Where:      "service = 'Microsoft.Insights' and action = 'activityLogs/read'",
			},
		},
		TableMap: map[string]*plugin.Table{
			"azure_alert_management":                                       tableAzureAlertMangement(ctx),
			"azure_api_management":                                         tableAzureAPIManagement(ctx),
			"azure_api_management_backend":                                 tableAzureAPIManagementBackend(ctx),
			"azure_app_configuration":                                      tableAzureAppConfiguration(ctx),
			"azure_app_service_environment":                                tableAzureAppServiceEnvironment(ctx),
			"azure_app_service_function_app":                               tableAzureAppServiceFunctionApp(ctx),
			"azure_app_service_plan":                                       tableAzureAppServicePlan(ctx),
			"azure_app_service_web_app":                                    tableAzureAppServiceWebApp(ctx),
			"azure_app_service_web_app_slot":                               tableAzureAppServiceWebAppSlot(ctx),
			"azure_application_gateway":                                    tableAzureApplicationGateway(ctx),
			"azure_application_insight":                                    tableAzureApplicationInsight(ctx),
			"azure_application_security_group":                             tableAzureApplicationSecurityGroup(ctx),
			"azure_automation_account":                                     tableAzureApAutomationAccount(ctx),
			"azure_automation_variable":                                    tableAzureApAutomationVariable(ctx),
			"azure_backup_policy":                                          tableAzureBackupPolicy(ctx),
			"azure_bastion_host":                                           tableAzureBastionHost(ctx),
			"azure_batch_account":                                          tableAzureBatchAccount(ctx),
			"azure_cdn_frontdoor_profile":                                  tableAzureCDNFrontDoorProfile(ctx),
			"azure_cognitive_account":                                      tableAzureCognitiveAccount(ctx),
			"azure_compute_availability_set":                               tableAzureComputeAvailabilitySet(ctx),
			"azure_compute_disk":                                           tableAzureComputeDisk(ctx),
			"azure_compute_disk_access":                                    tableAzureComputeDiskAccess(ctx),
			"azure_compute_disk_encryption_set":                            tableAzureComputeDiskEncryptionSet(ctx),
			"azure_compute_disk_metric_read_ops":                           tableAzureComputeDiskMetricReadOps(ctx),
			"azure_compute_disk_metric_read_ops_daily":                     tableAzureComputeDiskMetricReadOpsDaily(ctx),
			"azure_compute_disk_metric_read_ops_hourly":                    tableAzureComputeDiskMetricReadOpsHourly(ctx),
			"azure_compute_disk_metric_write_ops":                          tableAzureComputeDiskMetricWriteOps(ctx),
			"azure_compute_disk_metric_write_ops_daily":                    tableAzureComputeDiskMetricWriteOpsDaily(ctx),
			"azure_compute_disk_metric_write_ops_hourly":                   tableAzureComputeDiskMetricWriteOpsHourly(ctx),
			"azure_compute_image":                                          tableAzureComputeImage(ctx),
			"azure_compute_resource_sku":                                   tableAzureResourceSku(ctx),
			"azure_compute_snapshot":                                       tableAzureComputeSnapshot(ctx),
			"azure_compute_ssh_key":                                        tableAzureComputeSshKey(ctx),
			"azure_compute_virtual_machine":                                tableAzureComputeVirtualMachine(ctx),
			"azure_compute_virtual_machine_metric_available_memory":        tableAzureComputeVirtualMachineMetricAvailableMemory(ctx),
			"azure_compute_virtual_machine_metric_available_memory_daily":  tableAzureComputeVirtualMachineMetricAvailableMemoryDaily(ctx),
			"azure_compute_virtual_machine_metric_available_memory_hourly": tableAzureComputeVirtualMachineMetricAvailableMemoryHourly(ctx),
			"azure_compute_virtual_machine_metric_cpu_utilization":         tableAzureComputeVirtualMachineMetricCpuUtilization(ctx),
			"azure_compute_virtual_machine_metric_cpu_utilization_daily":   tableAzureComputeVirtualMachineMetricCpuUtilizationDaily(ctx),
			"azure_compute_virtual_machine_metric_cpu_utilization_hourly":  tableAzureComputeVirtualMachineMetricCpuUtilizationHourly(ctx),
			"azure_compute_virtual_machine_scale_set":                      tableAzureComputeVirtualMachineScaleSet(ctx),
			"azure_compute_virtual_machine_scale_set_network_interface":    tableAzureComputeVirtualMachineScaleSetNetworkInterface(ctx),
			"azure_compute_virtual_machine_scale_set_vm":                   tableAzureComputeVirtualMachineScaleSetVm(ctx),
			"azure_compute_virtual_machine_size":                           tableAzureComputeVirtualMachineSize(ctx),
			"azure_consumption_usage":                                      tableAzureConsumptionUsage(ctx),
			"azure_container_group":                                        tableAzureContainerGroup(ctx),
			"azure_container_registry":                                     tableAzureContainerRegistry(ctx),
			"azure_cosmosdb_account":                                       tableAzureCosmosDBAccount(ctx),
			"azure_cosmosdb_mongo_collection":                              tableAzureCosmosDBMongoCollection(ctx),
			"azure_cosmosdb_mongo_database":                                tableAzureCosmosDBMongoDatabase(ctx),
			"azure_cosmosdb_restorable_database_account":                   tableAzureCosmosDBRestorableDatabaseAccount(ctx),
			"azure_cosmosdb_sql_database":                                  tableAzureCosmosDBSQLDatabase(ctx),
			"azure_cost_by_resource_group_daily":                           tableAzureCostByResourceGroupDaily(ctx),
			"azure_cost_by_resource_group_monthly":                         tableAzureCostByResourceGroupMonthly(ctx),
			"azure_cost_by_service_daily":                                  tableAzureCostByServiceDaily(ctx),
			"azure_cost_by_service_monthly":                                tableAzureCostByServiceMonthly(ctx),
			"azure_cost_usage":                                             tableAzureCostUsage(ctx),
			"azure_data_factory":                                           tableAzureDataFactory(ctx),
			"azure_data_factory_dataset":                                   tableAzureDataFactoryDataset(ctx),
			"azure_data_factory_pipeline":                                  tableAzureDataFactoryPipeline(ctx),
			"azure_data_lake_analytics_account":                            tableAzureDataLakeAnalyticsAccount(ctx),
			"azure_data_lake_store":                                        tableAzureDataLakeStore(ctx),
			"azure_data_protection_backup_job":                             tableAzureDataProtectionBackupJob(ctx),
			"azure_data_protection_backup_vault":                           tableAzureDataProtectionBackupVault(ctx),
			"azure_databox_edge_device":                                    tableAzureDataBoxEdgeDevice(ctx),
			"azure_databricks_workspace":                                   tableAzureDatabricksWorkspace(ctx),
			"azure_diagnostic_setting":                                     tableAzureDiagnosticSetting(ctx),
			"azure_dns_zone":                                               tableAzureDNSZone(ctx),
			"azure_eventgrid_domain":                                       tableAzureEventGridDomain(ctx),
			"azure_eventgrid_topic":                                        tableAzureEventGridTopic(ctx),
			"azure_eventhub_namespace":                                     tableAzureEventHubNamespace(ctx),
			"azure_express_route_circuit":                                  tableAzureExpressRouteCircuit(ctx),
			"azure_firewall":                                               tableAzureFirewall(ctx),
			"azure_firewall_policy":                                        tableAzureFirewallPolicy(ctx),
			"azure_frontdoor":                                              tableAzureFrontDoor(ctx),
			"azure_hdinsight_cluster":                                      tableAzureHDInsightCluster(ctx),
			"azure_healthcare_service":                                     tableAzureHealthcareService(ctx),
			"azure_hpc_cache":                                              tableAzureHPCCache(ctx),
			"azure_hybrid_compute_machine":                                 tableAzureHybridComputeMachine(ctx),
			"azure_hybrid_kubernetes_connected_cluster":                    tableAzureHybridKubernetesConnectedCluster(ctx),
			"azure_iothub":                                                 tableAzureIotHub(ctx),
			"azure_iothub_dps":                                             tableAzureIotHubDps(ctx),
			"azure_key_vault":                                              tableAzureKeyVault(ctx),
			"azure_key_vault_certificate":                                  tableAzureKeyVaultCertificate(ctx),
			"azure_key_vault_deleted_vault":                                tableAzureKeyVaultDeletedVault(ctx),
			"azure_key_vault_key":                                          tableAzureKeyVaultKey(ctx),
			"azure_key_vault_key_version":                                  tableAzureKeyVaultKeyVersion(ctx),
			"azure_key_vault_managed_hardware_security_module":             tableAzureKeyVaultManagedHardwareSecurityModule(ctx),
			"azure_key_vault_secret":                                       tableAzureKeyVaultSecret(ctx),
			"azure_kubernetes_cluster":                                     tableAzureKubernetesCluster(ctx),
			"azure_kubernetes_service_version":                             tableAzureAKSVersion(ctx),
			"azure_kusto_cluster":                                          tableAzureKustoCluster(ctx),
			"azure_lb":                                                     tableAzureLoadBalancer(ctx),
			"azure_lb_backend_address_pool":                                tableAzureLoadBalancerBackendAddressPool(ctx),
			"azure_lb_nat_rule":                                            tableAzureLoadBalancerNatRule(ctx),
			"azure_lb_outbound_rule":                                       tableAzureLoadBalancerOutboundRule(ctx),
			"azure_lb_probe":                                               tableAzureLoadBalancerProbe(ctx),
			"azure_lb_rule":                                                tableAzureLoadBalancerRule(ctx),
			"azure_lighthouse_assignment":                                  tableAzureLighthouseAssignment(ctx),
			"azure_lighthouse_definition":                                  tableAzureLighthouseDefinition(ctx),
			"azure_location":                                               tableAzureLocation(ctx),
			"azure_log_alert":                                              tableAzureLogAlert(ctx),
			"azure_log_analytics_workspace":                                tableAzureLogAnalyticsWorkspace(ctx),
			"azure_log_profile":                                            tableAzureLogProfile(ctx),
			"azure_logic_app_workflow":                                     tableAzureLogicAppWorkflow(ctx),
			"azure_machine_learning_workspace":                             tableAzureMachineLearningWorkspace(ctx),
			"azure_maintenance_configuration":                              tableAzureMaintenanceConfiguration(ctx),
			"azure_management_group":                                       tableAzureManagementGroup(ctx),
			"azure_management_lock":                                        tableAzureManagementLock(ctx),
			"azure_mariadb_server":                                         tableAzureMariaDBServer(ctx),
			"azure_monitor_activity_log_event":                             tableAzureMonitorActivityLogEvent(ctx),
			"azure_monitor_log_profile":                                    tableAzureMonitorLogProfile(ctx),
			"azure_mssql_elasticpool":                                      tableAzureMSSQLElasticPool(ctx),
			"azure_mssql_managed_instance":                                 tableAzureMSSQLManagedInstance(ctx),
			"azure_mssql_virtual_machine":                                  tableAzureMSSQLVirtualMachine(ctx),
			"azure_mysql_flexible_server":                                  tableAzureMySQLFlexibleServer(ctx),
			"azure_mysql_server":                                           tableAzureMySQLServer(ctx),
			"azure_nat_gateway":                                            tableAzureNatGateway(ctx),
			"azure_network_interface":                                      tableAzureNetworkInterface(ctx),
			"azure_network_profile":                                        tableAzureNetworkProfile(ctx),
			"azure_network_security_group":                                 tableAzureNetworkSecurityGroup(ctx),
			"azure_network_watcher":                                        tableAzureNetworkWatcher(ctx),
			"azure_network_watcher_flow_log":                               tableAzureNetworkWatcherFlowLog(ctx),
			"azure_policy_assignment":                                      tableAzurePolicyAssignment(ctx),
			"azure_policy_definition":                                      tableAzurePolicyDefinition(ctx),
			"azure_postgresql_flexible_server":                             tableAzurePostgreSqlFlexibleServer(ctx),
			"azure_postgresql_server":                                      tableAzurePostgreSqlServer(ctx),
			"azure_private_dns_zone":                                       tableAzurePrivateDNSZone(ctx),
			"azure_private_endpoint":                                       tableAzurePrivateEndpoint(ctx),
			"azure_provider":                                               tableAzureProvider(ctx),
			"azure_public_ip":                                              tableAzurePublicIP(ctx),
			"azure_recovery_services_backup_job":                           tableAzureRecoveryServicesBackupJob(ctx),
			"azure_recovery_services_vault":                                tableAzureRecoveryServicesVault(ctx),
			"azure_redis_cache":                                            tableAzureRedisCache(ctx),
			"azure_resource":                                               tableAzureResourceResource(ctx),
			"azure_resource_group":                                         tableAzureResourceGroup(ctx),
			"azure_resource_link":                                          tableAzureResourceLink(ctx),
			"azure_role_assignment":                                        tableAzureIamRoleAssignment(ctx),
			"azure_role_definition":                                        tableAzureRoleDefinition(ctx),
			"azure_route_table":                                            tableAzureRouteTable(ctx),
			"azure_search_service":                                         tableAzureSearchService(ctx),
			"azure_security_center_auto_provisioning":                      tableAzureSecurityCenterAutoProvisioning(ctx),
			"azure_security_center_automation":                             tableAzureSecurityCenterAutomation(ctx),
			"azure_security_center_contact":                                tableAzureSecurityCenterContact(ctx),
			"azure_security_center_jit_network_access_policy":              tableAzureSecurityCenterJITNetworkAccessPolicy(ctx),
			"azure_security_center_setting":                                tableAzureSecurityCenterSetting(ctx),
			"azure_security_center_sub_assessment":                         tableAzureSecurityCenterSubAssessment(ctx),
			"azure_security_center_subscription_pricing":                   tableAzureSecurityCenterPricing(ctx),
			"azure_service_fabric_cluster":                                 tableAzureServiceFabricCluster(ctx),
			"azure_servicebus_namespace":                                   tableAzureServiceBusNamespace(ctx),
			"azure_signalr_service":                                        tableAzureSignalRService(ctx),
			"azure_spring_cloud_service":                                   tableAzureSpringCloudService(ctx),
			"azure_sql_database":                                           tableAzureSqlDatabase(ctx),
			"azure_sql_server":                                             tableAzureSQLServer(ctx),
			"azure_storage_account":                                        tableAzureStorageAccount(ctx),
			"azure_storage_blob":                                           tableAzureStorageBlob(ctx),
			"azure_storage_blob_service":                                   tableAzureStorageBlobService(ctx),
			"azure_storage_container":                                      tableAzureStorageContainer(ctx),
			"azure_storage_queue":                                          tableAzureStorageQueue(ctx),
			"azure_storage_share_file":                                     tableAzureStorageShareFile(ctx),
			"azure_storage_sync":                                           tableAzureStorageSync(ctx),
			"azure_storage_table":                                          tableAzureStorageTable(ctx),
			"azure_storage_table_service":                                  tableAzureStorageTableService(ctx),
			"azure_stream_analytics_job":                                   tableAzureStreamAnalyticsJob(ctx),
			"azure_subnet":                                                 tableAzureSubnet(ctx),
			"azure_subscription":                                           tableAzureSubscription(ctx),
			"azure_subscription_tenant_policy":                             tableAzureSubscriptionTenantPolicy(ctx),
			"azure_synapse_workspace":                                      tableAzureSynapseWorkspace(ctx),
			"azure_tenant":                                                 tableAzureTenant(ctx),
			"azure_virtual_network":                                        tableAzureVirtualNetwork(ctx),
			"azure_virtual_network_gateway":                                tableAzureVirtualNetworkGateway(ctx),
			"azure_web_application_firewall_policy":                        tableAzureWebApplicationFirewallPolicy(ctx),
		},
	}

	return p
}

func getSubscriptionIdForConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	subscriptionID, err := getSubscriptionIDMemoized(ctx, d, h)
	if err != nil {
		return nil, err
	}

	// The value must be returned as a string because connection-level quals do not support transform functions.
	// If the value is not returned as a string, queries that filter on subscription_id in the WHERE clause
	// (e.g., "SELECT id, subscription_id FROM azure_resource WHERE subscription_id = 'd46d7416...'")
	// will produce empty results due to a type mismatch during query evaluation.
	subscriptionIDStr := subscriptionID.(string)
	return subscriptionIDStr, nil
}
