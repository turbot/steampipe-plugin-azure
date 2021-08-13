select id, name, region
from azure.azure_search_service
where name = '{{ resourceName }}';
