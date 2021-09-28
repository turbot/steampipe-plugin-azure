select name, akas, title, tags
from azure.azure_spring_cloud_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
