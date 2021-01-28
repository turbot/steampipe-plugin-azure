select name, id
from azure.azure_storage_blob_service
where resource_group = '{{resourceName}}' and storage_account_name = 'dummy-{{resourceName}}'