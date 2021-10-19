select name, id, type, region, resource_group, subscription_id
from azure.azure_iothub_dps
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';