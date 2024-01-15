select name, akas, title
from azure.azure_spring_cloud_app
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
