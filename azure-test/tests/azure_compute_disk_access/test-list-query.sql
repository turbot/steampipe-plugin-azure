select name, id, type
from azure.azure_compute_disk_access
where name = '{{ resourceName }}';