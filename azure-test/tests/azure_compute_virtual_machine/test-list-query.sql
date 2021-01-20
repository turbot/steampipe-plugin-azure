select name, id, type
from azure.azure_compute_virtual_machine
where name = '{{resourceName}}'