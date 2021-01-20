select name, akas, tags, title
from azure.azure_cosmosdb_account
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
