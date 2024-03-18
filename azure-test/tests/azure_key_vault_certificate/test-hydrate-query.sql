select name, vault_name, id
from azure.azure_key_vault_certificate
where name = '{{resourceName}}' and title = '{{resourceName}}'
