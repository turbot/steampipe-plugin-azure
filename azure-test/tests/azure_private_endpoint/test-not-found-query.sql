select name, akas, tags, title
from azure.azure_private_endpoint
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
