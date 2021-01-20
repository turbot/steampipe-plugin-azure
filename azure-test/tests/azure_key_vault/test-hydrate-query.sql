select name, access_policies, akas, tags, title
from azure.azure_key_vault
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
