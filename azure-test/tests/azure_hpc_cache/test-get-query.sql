select name, id, type, region, resource_group, subscription_id, sku_name
from azure.azure_hpc_cache
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
