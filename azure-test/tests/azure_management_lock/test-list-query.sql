select id, name
from azure.azure_management_lock
where id = '{{ output.resource_id.value }}'
