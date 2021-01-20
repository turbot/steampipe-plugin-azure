select name, id, region, account_name, type, database_id, resource_group
from azure.azure_cosmosdb_mongo_database
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and account_name = '{{resourceName}}'
