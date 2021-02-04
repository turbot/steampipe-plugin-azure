select title, akas
from azure.azure_storage_blob
where resource_group = '{{resourceName}}' and storage_account_name = '{{resourceName}}' and container_name = '{{resourceName}}' and name = '{{resourceName}}'