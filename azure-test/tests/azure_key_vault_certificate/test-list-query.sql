select id, name
from azure.azure_key_vault_certificate
where name = '{{resourceName}}'
