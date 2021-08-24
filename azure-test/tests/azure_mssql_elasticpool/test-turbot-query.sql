select name, title, akas
from azure.azure_mssql_elasticpool
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';