select name, id, region, resource_group, subscription_id
from azure.azure_spring_cloud_app
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}' and service_name = '{{ resourceName }}';
