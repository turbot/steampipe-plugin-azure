select name, id, region, type, ip_configurations, sku_tier, resource_group
from azure.azure_firewall
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
