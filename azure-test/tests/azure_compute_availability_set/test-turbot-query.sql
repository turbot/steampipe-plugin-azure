select name, tags, title, akas
from azure.azure_compute_availability_set
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'