select name, akas, title
from azure_mysql_flexible_server
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
