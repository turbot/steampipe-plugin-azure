select name, akas, title
from azure.azure_subnet
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and virtual_network_name = '{{resourceName}}'
