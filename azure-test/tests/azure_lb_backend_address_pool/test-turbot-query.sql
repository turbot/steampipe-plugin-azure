select name, akas, title
from azure.azure_lb_backend_address_pool
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
