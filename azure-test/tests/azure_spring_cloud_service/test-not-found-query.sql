select name, id, type, region
from azure.azure_spring_cloud_service
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
