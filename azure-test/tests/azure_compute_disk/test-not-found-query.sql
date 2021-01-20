select name, tags, title, akas
from azure.azure_compute_disk
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'