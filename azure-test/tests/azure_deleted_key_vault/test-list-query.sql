select name, id, type
from azure.azure_deleted_key_vault
where name = '{{ output.resource_name.value }}';