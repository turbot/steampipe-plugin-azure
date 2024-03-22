select name, akas, title, tags
from azure_maintenance_configuration
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
