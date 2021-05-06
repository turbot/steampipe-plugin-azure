select name, akas, title
from azure.azure_security_center_setting
where name = '{{ output.resource_name.value }}';
