select name, id, type, region, resource_group, subscription_id
from azure.azure_signalr_service
where id = '{{ output.resource_id.value }}';
