select id, name, type
from azure.azure_policy_assignment
where name = '{{ output.resource_name.value }}'
