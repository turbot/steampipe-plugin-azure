select name, id
from azure.azure_storage_queue
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}' and storage_account_name = '{{resourceName}}'