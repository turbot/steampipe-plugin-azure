select name, id, type, region, resource_group
from azure.azure_network_watcher_flow_log
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}' and network_watcher_name = '{{ resourceName }}';
