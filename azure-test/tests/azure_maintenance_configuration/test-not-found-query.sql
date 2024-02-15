select name, akas, title
from azure_maintenance_configuration
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
