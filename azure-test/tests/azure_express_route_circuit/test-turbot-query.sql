select name, akas, title, tags
from azure.azure_express_route_circuit
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}'
