select name, title, akas
from azure.azure_mssql_managed_instance
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';