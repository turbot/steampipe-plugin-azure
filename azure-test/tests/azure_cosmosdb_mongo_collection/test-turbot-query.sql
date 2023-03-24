select 
  name,
  akas,
  title
from 
  azure.azure_cosmosdb_mongo_collection
where
  name = '{{resourceName}}'
  and resource_group = '{{resourceName}}'
  and database_name = '{{resourceName}}'
  and account_name = '{{resourceName}}';
