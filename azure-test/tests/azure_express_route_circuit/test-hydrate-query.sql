select name, sku_name, sku_tier, akas, tags, title
from azure.azure_express_route_circuit
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
