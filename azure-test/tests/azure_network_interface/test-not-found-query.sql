select name, akas, tags, title
from azure.azure_network_interface
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
