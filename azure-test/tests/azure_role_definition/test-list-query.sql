select id, name
from azure.azure_role_definition
where name = '{{ output.resource_id.value.split("/").pop() }}'
