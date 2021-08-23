select name, id, type
from azure.azure_mssql_managed_instance
where name = '{{ resourceName }}';