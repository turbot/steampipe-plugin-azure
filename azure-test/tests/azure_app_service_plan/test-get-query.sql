select name, id, kind, region, type, hyper_v, is_spot, is_xenon, reserved, per_site_scaling, resource_group, sku_tier, sku_size, sku_capacity
from azure.azure_app_service_plan
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'