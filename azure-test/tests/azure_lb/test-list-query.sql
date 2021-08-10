select name, id
from azure.azure_lb
where name = '{{ resourceName }}'
