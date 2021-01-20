select id, akas, title
from azure.azure_role_assignment
where id = '{{ output.resource_id.value }}'
