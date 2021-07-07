select id, name
from azure.azure_virtual_network_gateway
where name = '{{ resourceName }}';