select name, id, type, status, provisioning_state, sku_name, resource_group, region, subscription_id
from azure.azure_search_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
