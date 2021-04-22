select name, id, server_name, transparent_data_encryption
from azure.azure_sql_database
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';