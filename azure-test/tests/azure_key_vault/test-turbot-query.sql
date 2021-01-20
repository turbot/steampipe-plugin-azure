select name, akas, title, tags
from azure.azure_key_vault
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
