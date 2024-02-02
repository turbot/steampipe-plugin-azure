select name, id, type, region
from azure.azure_spring_cloud_app
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
