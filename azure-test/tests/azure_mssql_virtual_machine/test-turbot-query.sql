select name, title
from azure.azure_mssql_virtual_machine
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
