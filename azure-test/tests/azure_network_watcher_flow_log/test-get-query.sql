select name, id, enabled, network_watcher_name, type, file_type, retention_policy_days, retention_policy_enabled, version, storage_id, target_resource_id, traffic_analytics, region, resource_group, subscription_id
from azure.azure_network_watcher_flow_log
where name = 'Microsoft.Network{{ resourceName }}{{ resourceName }}' and resource_group = '{{ resourceName }}' and network_watcher_name = '{{ resourceName }}';
