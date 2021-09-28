select name, id, type, region
from azure.azure_application_gateway
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
