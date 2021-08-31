select name, id, type, resource_group, subscription_id
from azure.azure_lb_backend_address_pool
where load_balancer_name = '{{ resourceName }}' and name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
