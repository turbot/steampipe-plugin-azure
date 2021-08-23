select name, title, akas
from azure.azure_mssql_managed_instance
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';