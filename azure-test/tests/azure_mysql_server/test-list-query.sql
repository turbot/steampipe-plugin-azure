select id, name, location
from azure.azure_mysql_server
where name = '{{ resourceName }}';
