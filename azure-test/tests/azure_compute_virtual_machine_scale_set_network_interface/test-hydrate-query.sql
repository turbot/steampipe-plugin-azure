select name, id
from azure.azure_compute_virtual_machine_scale_set_network_interface
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and id = '{{ output.resource_id.value }}';