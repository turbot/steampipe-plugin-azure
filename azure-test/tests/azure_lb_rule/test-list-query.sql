select name, id
from azure.azure_lb_rule
where name = '{{ resourceName }}';
