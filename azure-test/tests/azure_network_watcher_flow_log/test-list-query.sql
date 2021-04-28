select id, name
from azure.azure_network_watcher_flow_log
where name = 'Microsoft.Network{{ resourceName }}{{ resourceName }}';
