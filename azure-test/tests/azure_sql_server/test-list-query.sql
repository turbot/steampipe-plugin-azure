select id, name, location
from azure.azure_sql_server
where name = '{{ resourceName }}';
