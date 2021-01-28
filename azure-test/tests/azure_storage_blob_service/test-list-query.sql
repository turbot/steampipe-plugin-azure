select name, id
from azure.azure_storage_blob_service
where id = '{{ output.resource_id.value }}'