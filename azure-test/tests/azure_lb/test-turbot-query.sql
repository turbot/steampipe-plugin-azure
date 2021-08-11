select name, akas, tags, title
from azure.azure_lb
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}'
