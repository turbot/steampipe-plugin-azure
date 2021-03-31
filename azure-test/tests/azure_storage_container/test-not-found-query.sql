select name, id
from azure.azure_storage_container
where name = 'dummy{{resourceName}}' and resource_group = '{{resourceName}}' and
account_name = '{{resourceName}}'