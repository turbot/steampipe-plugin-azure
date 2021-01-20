select name, id, role_name, type, role_type, description, assignable_scopes, permissions
from azure.azure_role_definition
where name = '{{ output.resource_id.value.split("/").pop() }}'
