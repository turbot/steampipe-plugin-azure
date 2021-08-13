select name, akas, title
from azure.azure_lb_outbound_rule
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
