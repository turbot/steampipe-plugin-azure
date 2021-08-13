select name, id
from azure.azure_lb_probe
where name = 'dummy-test-{{ resourceName }}' and resource_group = '{{ resourceName }}';
