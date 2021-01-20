select name, id, type, region, resource_group
from azure.azure_network_watcher
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
