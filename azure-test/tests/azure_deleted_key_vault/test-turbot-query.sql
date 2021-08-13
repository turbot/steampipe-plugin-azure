select name, title, akas
from azure.azure_deleted_key_vault
where name = '{{ output.resource_name.value }}';