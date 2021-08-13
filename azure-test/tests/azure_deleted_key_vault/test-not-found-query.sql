select name, tags, title, akas
from azure.azure_deleted_key_vault
where name = 'dummy-{{ output.resource_name.value }}';