select name, id, type
from azure.azure_iothub
where name = '{{ resourceName }}';