select name, akas, title, tags
from azure.azure_route_table
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
