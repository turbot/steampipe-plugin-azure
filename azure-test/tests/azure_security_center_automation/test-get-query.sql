select name, id
from azure_security_center_automation
where name = '{{ output.resource_name.value }}' and resource_group = '{{ output.resource_name.value }}';
