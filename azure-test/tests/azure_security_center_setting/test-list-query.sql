select id, name
from azure.azure_security_center_setting
where id = '{{ output.resource_id.value }}'
