select name, id
from azure.azure_key_vault_certificate
where name = '{{resourceName}}' and vault_name = '{{resourceName}}'
