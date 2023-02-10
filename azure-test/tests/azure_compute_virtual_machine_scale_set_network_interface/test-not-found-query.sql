select name, id
from azure_compute_virtual_machine_scale_set_network_interface
where name = 'dummy{{ resourceName }}' and resource_group = '{{ resourceName }}';