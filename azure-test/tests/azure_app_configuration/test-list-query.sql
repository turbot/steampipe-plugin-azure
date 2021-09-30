select name, id, type, region, resource_group, subscription_id
from azure.azure_app_configuration
where id = '{{ output.resource_id.value }}';
