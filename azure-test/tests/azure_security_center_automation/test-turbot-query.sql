select name, akas, title
from azure_security_center_automation
where name = '{{ output.resource_name.value }}';
