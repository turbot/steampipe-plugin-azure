select name, akas, title
from azure_postgresql_flexible_server
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
