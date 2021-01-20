select name, id, region, type, disable_bgp_route_propagation, routes, resource_group
from azure.azure_route_table
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
