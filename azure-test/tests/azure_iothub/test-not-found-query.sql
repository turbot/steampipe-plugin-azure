select name, tags, title, akas
from azure.azure_iothub
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';