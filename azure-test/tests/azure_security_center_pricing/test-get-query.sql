select name, id, type
from azure.azure_security_center_contact
where name = '{{ output.resource_name.value }}';
