select name, id, scope, type, principal_id, role_definition_id
from azure.azure_role_assignment_for_resource
where name = '{{ output.resource_id.value.split("/").pop() }}'