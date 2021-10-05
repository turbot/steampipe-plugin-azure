select name, id, type
from azure.azure_iothub_dps
where name = '{{ resourceName }}';