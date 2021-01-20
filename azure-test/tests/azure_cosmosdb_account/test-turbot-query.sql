select name, akas, title, tags
from azure.azure_cosmosdb_account
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
