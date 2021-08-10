select name, id
from azure.azure_virtual_machine_scale_set
where name = 'dummy{{ resourceName }}' and resource_group = '{{ resourceName }}';