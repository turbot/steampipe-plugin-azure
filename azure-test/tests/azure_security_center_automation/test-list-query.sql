select id, name
from azure_security_center_automation
where id = '{{ output.resource_id.value }}'
