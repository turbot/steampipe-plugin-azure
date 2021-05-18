select id, name, type
from azure.azure_policy_assignment
where id = '{{ output.resource_id.value }}'
