select name, akas, title, tags
from azure.azure_sql_database
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';