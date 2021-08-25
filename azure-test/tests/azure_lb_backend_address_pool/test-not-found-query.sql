select name, id
from azure.azure_lb_backend_address_pool
where name = 'dummy-test-{{ resourceName }}' and resource_group = '{{ resourceName }}';
