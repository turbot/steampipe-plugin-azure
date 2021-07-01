select name, id, region, type, enable_bgp, resource_group
from azure.azure_virtual_network_gateway
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
