select name, id, type, region
from azure.azure_lb
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
