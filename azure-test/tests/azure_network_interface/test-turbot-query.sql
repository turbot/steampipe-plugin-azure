select name, akas, title, tags
from azure.azure_network_interface
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
