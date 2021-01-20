select name, akas, tags, title
from azure.azure_public_ip
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
