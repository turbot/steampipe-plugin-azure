select name, akas, title
from azure.azure_app_configuration
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
