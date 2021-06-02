select id, name, display_name
from azure.azure_policy_definition
where display_name = '{{ output.resource_name.value }}';
