select name, id, type
from azure.azure_security_center_setting
where name = '{{ output.resource_name.value }}';
