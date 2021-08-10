select name, id, type
from azure.azure_virtual_machine_scale_set
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';