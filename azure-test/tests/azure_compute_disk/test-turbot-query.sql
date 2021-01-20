select name, tags, title, akas
from azure.azure_compute_disk
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'