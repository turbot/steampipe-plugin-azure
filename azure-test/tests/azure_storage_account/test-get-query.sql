select name, id, access_tier, sku_tier, sku_name, kind
from azure.azure_storage_account
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'