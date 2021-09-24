select name, id, type, region, resource_group, subscription_id, tags
from azure.azure_storage_sync
where id = '{{ output.resource_id.value }}';
