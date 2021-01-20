select name, sku_name
from azure.azure_api_management
where name = '{{resourceName}}' and resource_group = '{{ output.resource_group_name.value }}'
