select id, name
from azure.azure_security_center_pricing
where id = '{{ output.resource_id.value }}'
