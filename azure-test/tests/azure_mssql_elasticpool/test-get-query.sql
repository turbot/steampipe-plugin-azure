select name, id, type, resource_group, subscription_id
from azure.azure_mssql_elasticpool
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}' and server_name = '{{ resourceName }}';