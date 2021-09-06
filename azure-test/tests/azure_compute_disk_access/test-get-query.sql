select name, id, type, resource_group, subscription_id
from azure.azure_compute_disk_access
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';