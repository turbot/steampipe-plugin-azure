select name, akas, tags, title
from azure.azure_virtual_network
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
