select name, akas, title
from azure.azure_databox_edge_device
where name = 'dummy-{{ resourceName }}';
