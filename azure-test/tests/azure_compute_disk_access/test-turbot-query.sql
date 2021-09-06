select name, title, akas
from azure.azure_compute_disk_access
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';