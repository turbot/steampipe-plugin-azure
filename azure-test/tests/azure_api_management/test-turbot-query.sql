select name, akas, title
from azure.azure_api_management
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
