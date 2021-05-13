select name, akas, title
from azure.azure_security_center_contact
where name = '{{ output.resource_name.value }}';
