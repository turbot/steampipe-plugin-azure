select name, id, load_balancer_name, type, resource_group, subscription_id
from azure.azure_lb_outbound_rule
where name = '{{ resourceName }}' and load_balancer_name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
