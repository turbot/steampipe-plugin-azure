select name, id, region, resource_group
from azure.azure_maintenance_configuration
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
