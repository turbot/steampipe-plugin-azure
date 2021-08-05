select name, id, type, resource_group, subscription_id
from azure.azure_data_lake_store
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';