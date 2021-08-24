select name, id
from azure.azure_lb_outbound_rule
where name = '{{ resourceName }}';
