select name, akas, title
from azure.azure_policy_assignment
where name = '{{ output.resource_name.value }}';
