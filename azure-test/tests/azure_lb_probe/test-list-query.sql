select name, id
from azure.azure_lb_probe
where name = '{{ resourceName }}'
