select name, id
from azure_compute_virtual_machine_scale_set_network_interface
where id = '{{ output.resource_id.value }}';