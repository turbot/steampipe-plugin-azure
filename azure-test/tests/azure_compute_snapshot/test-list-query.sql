select name, id, type
from azure.azure_compute_snapshot
where name = '{{resourceName}}snapshot'