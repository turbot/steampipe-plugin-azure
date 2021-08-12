select name, id, type, resource_group, subscription_id
from azure.azure_data_lake_analytics_account
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';