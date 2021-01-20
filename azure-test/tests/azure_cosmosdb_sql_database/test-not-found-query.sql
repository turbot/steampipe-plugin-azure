select name, akas, tags, title
from azure.azure_cosmosdb_sql_database
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}' and account_name = '{{resourceName}}'
