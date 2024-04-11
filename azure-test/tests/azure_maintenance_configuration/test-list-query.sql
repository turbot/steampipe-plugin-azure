select id, name, region
from azure_maintenance_configuration
where name = '{{ resourceName }}';
