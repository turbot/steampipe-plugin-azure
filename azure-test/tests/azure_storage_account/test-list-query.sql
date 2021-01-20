select name, id
from azure.azure_storage_account
where id = '{{ output.resource_id.value }}'