select name, id
from azure.azure_storage_table
where id = '{{ output.resource_id.value }}'