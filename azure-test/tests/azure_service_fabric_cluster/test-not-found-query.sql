select name, id, type, region
from azure.azure_service_fabric_cluster
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
