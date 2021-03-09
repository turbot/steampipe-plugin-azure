select name, akas, title
from azure.azure_sql_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
