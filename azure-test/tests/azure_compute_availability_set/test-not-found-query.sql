select name, tags, title, akas
from azure.azure_compute_availability_set
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'