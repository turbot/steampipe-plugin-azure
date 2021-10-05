select name, title, akas
from azure.azure_iothub_dps
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';