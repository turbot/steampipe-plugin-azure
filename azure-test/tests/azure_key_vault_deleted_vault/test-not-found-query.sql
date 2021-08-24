select name, tags, title, akas
from azure.azure_key_vault_deleted_vault
where name = 'dummy-{{ output.resource_name.value }}';