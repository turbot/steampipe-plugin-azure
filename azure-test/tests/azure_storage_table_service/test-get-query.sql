select name, id, storage_account_name, type, cors_rules
from azure.azure_storage_table_service
where resource_group = '{{resourceName}}' and storage_account_name = '{{resourceName}}'