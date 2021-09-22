select name, id, type, resource_group, subscription_id
from azure.azure_compute_disk_access
where name = '{{ output.resource_name.value }}' and resource_group = '{{ output.resource_name.value }}';