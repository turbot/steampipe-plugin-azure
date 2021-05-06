select id, name
from azure.azure_security_center_auto_provisioning
where id = '{{ output.resource_id.value }}'
