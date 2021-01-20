select name, akas, tags, title
from azure.azure_resource_group
where name = '{{resourceName}}'
