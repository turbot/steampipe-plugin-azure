select name, id, type, region, resource_group, subscription_id, tags
from azure.azure_spring_cloud_service
where id = '{{ output.resource_id.value }}';
