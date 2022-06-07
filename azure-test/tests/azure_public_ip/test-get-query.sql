select name, id, region, type, resource_group
from azure.azure_public_ip
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
