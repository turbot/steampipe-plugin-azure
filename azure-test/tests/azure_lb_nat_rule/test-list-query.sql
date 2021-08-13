select name, id
from azure.azure_lb_nat_rule
where name = '{{ resourceName }}';
