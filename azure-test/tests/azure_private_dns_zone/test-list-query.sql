select id,
    name
from azure.azure_private_dns_zone
where id = '{{ output.resource_id.value }}'
