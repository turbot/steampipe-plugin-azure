select name, tags, title, akas
from azure.azure_compute_virtual_machine
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'