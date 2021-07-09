select name, akas, title
from azure.azure_virtual_network_gateway
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';