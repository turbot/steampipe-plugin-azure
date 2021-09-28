select name, akas, title
from azure.azure_application_gateway
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
