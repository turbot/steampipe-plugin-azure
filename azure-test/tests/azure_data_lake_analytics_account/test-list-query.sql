select name, id, type
from azure.azure_data_lake_analytics_account
where name = '{{ resourceName }}';