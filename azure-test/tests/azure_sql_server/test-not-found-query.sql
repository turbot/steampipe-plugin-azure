select name, akas, title
from azure.azure_sql_server
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
