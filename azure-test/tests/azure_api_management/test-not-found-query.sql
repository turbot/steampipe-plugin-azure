select name, id, region, type, publisher_email, publisher_name, sku_capacity
from azure.azure_api_management
where name = '{{resourceName}}' and resource_group = 'dummy-{{resourceName}}'
