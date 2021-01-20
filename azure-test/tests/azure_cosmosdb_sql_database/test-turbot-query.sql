select name, akas, title
from azure.azure_cosmosdb_sql_database
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and account_name = '{{resourceName}}'
