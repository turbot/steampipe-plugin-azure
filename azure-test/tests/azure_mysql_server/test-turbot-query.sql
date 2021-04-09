select name, akas, title, tags
from azure.azure_mysql_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
