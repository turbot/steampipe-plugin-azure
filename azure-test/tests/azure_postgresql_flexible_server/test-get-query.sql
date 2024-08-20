select name, version, administrator_login, resource_group, subscription_id
from azure_postgresql_flexible_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
