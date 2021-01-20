select name, akas, title, tags
from azure.azure_public_ip
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
