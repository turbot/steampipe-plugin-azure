select name, type, title
from azure.azure_policy_assignment
where name = 'dummy-{{ output.resource_name.value }}';
