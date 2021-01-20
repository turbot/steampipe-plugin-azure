select name, id, access_tier
from azure.azure_storage_account
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'