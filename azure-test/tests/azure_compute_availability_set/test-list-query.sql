select name, id, type
from azure.azure_compute_availability_set
where name = '{{resourceName}}'