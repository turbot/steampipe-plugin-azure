select name, akas, tags, title
from azure.azure_virtual_network
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
