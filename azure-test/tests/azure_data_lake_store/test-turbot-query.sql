select title, akas
from azure.azure_data_lake_store
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';