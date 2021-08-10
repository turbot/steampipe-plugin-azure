select name, id, sku_tier, sku_name, tags
from azure.azure_virtual_machine_scale_set
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';