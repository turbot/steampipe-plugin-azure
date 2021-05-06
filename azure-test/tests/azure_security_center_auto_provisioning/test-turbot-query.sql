select name, akas, title
from azure.azure_security_center_auto_provisioning
where name = '{{ output.resource_name.value }}';
