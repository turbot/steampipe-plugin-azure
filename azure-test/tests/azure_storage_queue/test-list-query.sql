select name, id
from azure.azure_storage_queue
where id = '{{ output.resource_id.value }}'