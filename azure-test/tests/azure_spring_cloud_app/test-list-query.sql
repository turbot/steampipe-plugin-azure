select name, id, region, resource_group, subscription_id
from azure.azure_spring_cloud_app
where id = '{{ output.resource_id.value }}';
