select name, replica_count, partition_count
from azure.azure_search_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';