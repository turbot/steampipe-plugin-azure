select name, id
from azure.azure_storage_table_service
where id = '{{ output.resource_id.value }}'