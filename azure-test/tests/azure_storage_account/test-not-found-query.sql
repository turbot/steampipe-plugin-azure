select name, id
from azure.azure_storage_account
where name = 'dummy{{resourceName}}' and resource_group = '{{resourceName}}'