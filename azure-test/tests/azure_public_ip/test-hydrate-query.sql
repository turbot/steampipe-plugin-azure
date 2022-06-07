select name, sku_name, public_ip_allocation_method, akas, tags, title
from azure.azure_public_ip
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
