select name, title, tags, akas
from azure.azure_compute_disk_encryption_set
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'