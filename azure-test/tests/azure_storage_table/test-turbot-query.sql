select title, akas
from azure.azure_storage_table
where resource_group = '{{resourceName}}' and storage_account_name = '{{resourceName}}' and name = '{{resourceName}}'