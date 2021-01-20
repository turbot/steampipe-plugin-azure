select name, id, scope, type, principal_id, principal_type, role_definition_id
from azure.azure_role_assignment
where id = '{{ output.resource_id.value }}'
