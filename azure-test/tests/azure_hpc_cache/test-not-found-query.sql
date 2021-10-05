select name, id, type, region
from azure.azure_hpc_cache
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
