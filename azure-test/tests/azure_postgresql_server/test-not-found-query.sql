select name, akas, title
from azure.azure_postgresql_server
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
