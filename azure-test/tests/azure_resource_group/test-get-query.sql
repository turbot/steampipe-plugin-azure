select name, id, region, type
from azure.azure_resource_group
where name = '{{resourceName}}'
