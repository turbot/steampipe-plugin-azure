select name, akas, tags, title
from azure.azure_cosmosdb_account
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
