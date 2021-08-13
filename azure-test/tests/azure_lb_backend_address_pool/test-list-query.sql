select name, id
from azure.azure_lb_backend_address_pool
where name = '{{ resourceName }}';
