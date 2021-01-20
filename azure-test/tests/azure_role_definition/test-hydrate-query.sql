select name, akas, title
from azure.azure_role_definition
where name = '{{ output.resource_id.value.split("/").pop() }}'
