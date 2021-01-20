select name, tags, title, akas
from azure.azure_compute_snapshot
where name = '{{resourceName}}snapshot' and resource_group = '{{resourceName}}'