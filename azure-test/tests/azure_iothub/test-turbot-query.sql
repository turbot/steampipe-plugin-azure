select name, title, akas
from azure.azure_iothub
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';