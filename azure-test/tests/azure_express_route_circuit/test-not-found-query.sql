select name, akas, tags, title
from azure.azure_express_route_circuit
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';