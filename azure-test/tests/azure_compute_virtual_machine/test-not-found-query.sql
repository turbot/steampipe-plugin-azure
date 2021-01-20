select name, tags, title, akas
from azure.azure_compute_virtual_machine
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'