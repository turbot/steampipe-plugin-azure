select name, id, type, region
from azure.azure_app_configuration
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
