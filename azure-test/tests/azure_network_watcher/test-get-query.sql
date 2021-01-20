select name, id, type, region, resource_group
from azure.azure_network_watcher
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
