select name, akas, title
from azure.azure_search_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
