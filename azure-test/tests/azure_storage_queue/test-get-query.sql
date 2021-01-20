select name, id, storage_account_name, type, metadata
from azure.azure_storage_queue
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and storage_account_name = '{{resourceName}}'