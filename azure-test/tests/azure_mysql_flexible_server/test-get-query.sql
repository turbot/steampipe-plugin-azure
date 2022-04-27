select name, version, administrator_login, fully_qualified_domain_name, resource_group, subscription_id
from azure_mysql_flexible_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
