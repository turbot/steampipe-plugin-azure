select name, id, type
from azure.azure_key_vault_deleted_vault
where name = '{{ output.resource_name.value }}';