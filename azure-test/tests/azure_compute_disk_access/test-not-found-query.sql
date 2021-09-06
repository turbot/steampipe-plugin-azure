select name, tags, title, akas
from azure.azure_compute_disk_access
where name = 'dummy-{{ output.resource_name.value }}' and resource_group = '{{ output.resource_name.value }}';