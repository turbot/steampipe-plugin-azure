# Azure Plugin Rate Limiter Tags

This document lists the rate limiter tags for each table in the Azure plugin. These tags are used to control API request rates and ensure compliance with Azure service limits.

## Table Rate Limiter Tags

| Table Name | Service | Action |
|------------|---------|---------|
| azure_alert_management | Microsoft.AlertsManagement | alerts/read |
| azure_api_management | Microsoft.ApiManagement | service/read |
| azure_app_configuration | Microsoft.AppConfiguration | configurationStores/read |
| azure_batch_account | Microsoft.Batch | batchAccounts/read |
| azure_cognitive_account | Microsoft.CognitiveServices | accounts/read |
| azure_cosmosdb_account | Microsoft.DocumentDB | databaseAccounts/read |
| azure_cosmosdb_restorable_database_account | Microsoft.DocumentDB | restorableDatabaseAccounts/read |
| azure_cosmosdb_sql_database | Microsoft.DocumentDB | sqlDatabases/read |
| azure_data_lake_analytics_account | Microsoft.DataLakeAnalytics | accounts/read |
| azure_data_lake_store | Microsoft.DataLakeStore | accounts/read |
| azure_frontdoor | Microsoft.Network | frontDoors/read |
| azure_hpc_cache | Microsoft.StorageCache | caches/read |
| azure_hybrid_kubernetes_connected_cluster | Microsoft.Kubernetes | connectedClusters/read |
| azure_lb | Microsoft.Network | loadBalancers/read |
| azure_lb_nat_rule | Microsoft.Network | loadBalancers/inboundNatRules/read |
| azure_lb_outbound_rule | Microsoft.Network | loadBalancers/outboundRules/read |
| azure_lb_rule | Microsoft.Network | loadBalancers/loadBalancingRules/read |
| azure_log_alert | Microsoft.Insights | activityLogAlerts/read |
| azure_log_analytics_workspace | Microsoft.OperationalInsights | workspaces/read |
| azure_maintenance_configuration | Microsoft.Maintenance | maintenanceConfigurations/read |
| azure_management_lock | Microsoft.Authorization | locks/read |
| azure_monitor_activity_log_event | Microsoft.Insights | activityLogs/read |
| azure_monitor_log_profile | Microsoft.Insights | logProfiles/read |
| azure_mssql_elasticpool | Microsoft.Sql | elasticPools/read |
| azure_mysql_flexible_server | Microsoft.DBforMySQL | flexibleServers/read |
| azure_network_profile | Microsoft.Network | networkProfiles/read |
| azure_network_watcher | Microsoft.Network | networkWatchers/read |
| azure_postgresql_server | Microsoft.DBforPostgreSQL | servers/read |
| azure_redis_cache | Microsoft.Cache | redis/read |
| azure_resource | Microsoft.Resources | resources/read |
| azure_route_table | Microsoft.Network | routeTables/read |
| azure_security_center_jit_network_access_policy | Microsoft.Security | jitNetworkAccessPolicies/read |
| azure_servicebus_namespace | Microsoft.ServiceBus | namespaces/read |
| azure_sql_database | Microsoft.Sql | databases/read |
| azure_sql_server | Microsoft.Sql | servers/read |
| azure_storage_account | Microsoft.Storage | storageAccounts/read |
| azure_storage_blob | Microsoft.Storage | blobs/read |
| azure_storage_table | Microsoft.Storage | storageAccounts/tableServices/tables/read |
| azure_storage_table_service | Microsoft.Storage | storageAccounts/tableServices/read |

## Usage

These tags are used in both Get and List operations for each table. They help:

1. Identify which Azure service is being accessed
2. Specify the exact action being performed
3. Control rate limiting for API requests
4. Track API usage and quotas

The format is:
```go
Tags: map[string]string{
    "service": "Microsoft.ServiceName",
    "action":  "resource/action",
}
```

## Rate Limiting Implementation

The plugin implements rate limiting through:

1. Automatic retry mechanisms using `ApplyRetryRules`
2. Wait periods between paginated requests using `WaitForListRateLimit`
3. Service-specific rate limits based on Azure API quotas

This ensures the plugin operates within Azure service limits while maintaining reliable data collection. 