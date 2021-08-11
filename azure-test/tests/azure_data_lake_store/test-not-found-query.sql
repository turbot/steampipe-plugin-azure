select name
from azure.azure_data_lake_store
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';