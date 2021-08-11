select name, id, type
from azure.azure_data_lake_store
where name = '{{ resourceName }}';