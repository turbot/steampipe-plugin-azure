select 
  name,
  id,
  region,
  account_name,
  database_name,
  type,
  resource_group
from 
  azure.azure_cosmosdb_mongo_collection
where
  name = '{{resourceName}}-dummy'
  and resource_group = '{{resourceName}}'
  and database_name = '{{resourceName}}'
  and account_name = '{{resourceName}}';
