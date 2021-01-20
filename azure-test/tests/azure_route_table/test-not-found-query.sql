select name, akas, tags, title
from azure.azure_route_table
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
