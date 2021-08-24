select name, id
from azure.azure_lb_nat_rule
where name = 'dummy-test-{{ resourceName }}' and resource_group = '{{ resourceName }}';
