select id, name, location
from azure.azure_postgresql_server
where name = '{{ resourceName }}';
