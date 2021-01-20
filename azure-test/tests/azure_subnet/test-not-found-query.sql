select name, akas, title
from azure.azure_subnet
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}' and virtual_network_name = '{{resourceName}}'
