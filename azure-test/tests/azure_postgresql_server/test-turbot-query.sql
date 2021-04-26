select name, akas, title, tags
from azure.azure_postgresql_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
