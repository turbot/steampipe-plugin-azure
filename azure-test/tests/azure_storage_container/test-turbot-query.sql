select title, akas
from azure.azure_storage_container
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and
account_name = '{{resourceName}}'