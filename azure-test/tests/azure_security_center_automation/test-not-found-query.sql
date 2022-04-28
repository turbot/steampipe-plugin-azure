select name, akas, title
from azure_security_center_automation
where name = 'dummy-{{ output.resource_name.value }}';
