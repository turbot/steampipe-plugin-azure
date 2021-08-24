select name, id, type
from azure.azure_mssql_elasticpool
where name = '{{ resourceName }}';