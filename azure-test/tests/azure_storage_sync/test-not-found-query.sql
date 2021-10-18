select name, id, type, region
from azure.azure_storage_sync
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
