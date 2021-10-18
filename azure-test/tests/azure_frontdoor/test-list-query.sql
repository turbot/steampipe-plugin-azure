select name, id, type, resource_group, subscription_id
from azure.azure_frontdoor
where id = '{{ output.resource_id.value }}';
