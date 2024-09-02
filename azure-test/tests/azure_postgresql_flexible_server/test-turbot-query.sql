select name, akas, title, tags
from azure_postgresql_flexible_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
