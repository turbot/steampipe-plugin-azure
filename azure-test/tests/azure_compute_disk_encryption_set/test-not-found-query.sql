select name, id, type
from azure.azure_compute_disk_encryption_set
where name = 'dummy{{resourceName}}' and resource_group = '{{resourceName}}'