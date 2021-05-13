select id, name
from azure.azure_security_center_contact
where id = '{{ output.resource_id.value }}'
