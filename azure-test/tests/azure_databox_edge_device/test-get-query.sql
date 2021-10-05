select name, id,type
from azure.azure_databox_edge_device
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
