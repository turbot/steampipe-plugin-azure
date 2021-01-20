select name, akas, tags, title
from azure.azure_network_watcher
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
