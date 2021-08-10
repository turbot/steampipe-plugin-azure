select name, id 
from azure.azure_virtual_machine_scale_set
where id = '{{ output.resource_id.value }}';