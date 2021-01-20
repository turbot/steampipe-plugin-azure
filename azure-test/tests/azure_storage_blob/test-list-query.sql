select name, id
from azure.azure_storage_blob
where id = '{{ output.resource_id.value }}'