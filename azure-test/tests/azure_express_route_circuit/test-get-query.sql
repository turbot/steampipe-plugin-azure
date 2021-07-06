select name, id, sku_tier, resource_group
from azure.azure_express_route_circuit
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}'
