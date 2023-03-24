select 
  name,
  id,
  region,
  account_name,
  database_name,
  resource_group
from 
  azure.azure_cosmosdb_mongo_collection
where
  name = '{{resourceName}}'
  and database_name = '{{resourceName}}';
