select name, id, enabled, network_watcher_name, type, flowlog_format_type, flowlog_retention_days, flowlog_retention_enabled, flowlog_version, storage_id, target_resource_id, flow_log_traffic_analytics_configuration, region, resource_group, subscription_id
from azure.azure_network_watcher_flow_log
where name = 'Microsoft.Network{{ resourceName }}{{ resourceName }}' and resource_group = '{{ resourceName }}' and network_watcher_name = '{{ resourceName }}';
