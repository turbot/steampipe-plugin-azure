select name, id, type
from azure.azure_security_center_auto_provisioning
where name = '{{ output.resource_name.value }}';
