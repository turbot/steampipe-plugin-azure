select name, akas, title
from azure.azure_security_center_setting
where name = 'dummy-{{ output.resource_name.value }}';
