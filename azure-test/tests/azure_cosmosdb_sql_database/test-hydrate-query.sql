select name, id, region, account_name, type, database_colls, database_id, database_users, resource_group
from azure.azure_cosmosdb_sql_database
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and account_name = '{{resourceName}}'
