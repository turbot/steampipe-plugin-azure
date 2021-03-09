select name, type, kind, location, administrator_login, version, fully_qualified_domain_name, tags_src, resource_group, region, subscription_id
from azure.azure_sql_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
