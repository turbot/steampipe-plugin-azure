select name, id, type, resource_group, subscription_id
from azure.azure_mssql_managed_instance
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';