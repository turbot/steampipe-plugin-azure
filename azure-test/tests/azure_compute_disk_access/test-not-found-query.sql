select name, tags, title, akas
from azure.azure_compute_disk_access
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';