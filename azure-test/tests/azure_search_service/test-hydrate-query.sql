select name, replica_count, partition_count, network_rule_set
from azure.azure_search_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';