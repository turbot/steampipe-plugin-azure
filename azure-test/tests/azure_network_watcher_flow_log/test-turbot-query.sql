select name, akas, title
from azure.azure_network_watcher_flow_log
where name = 'Microsoft.Network{{ resourceName }}{{ resourceName }}' and resource_group = '{{ resourceName }}' and network_watcher_name = '{{ resourceName }}';
