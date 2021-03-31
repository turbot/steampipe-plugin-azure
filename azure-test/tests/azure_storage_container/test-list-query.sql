select name, id
from azure.azure_storage_container
where id = '{{ output.resource_id.value }}'