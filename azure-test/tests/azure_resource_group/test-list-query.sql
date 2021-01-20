select id, name
from azure.azure_resource_group
where id = '{{ output.resource_id.value }}'
