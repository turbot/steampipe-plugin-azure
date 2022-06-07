select id, name
from azure.azure_public_ip
where name = '{{ resourceName }}';
