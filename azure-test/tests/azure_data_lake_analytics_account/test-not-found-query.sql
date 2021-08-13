select name
from azure.azure_data_lake_analytics_account
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';