select name, akas, tags, title
from azure.azure_key_vault
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
