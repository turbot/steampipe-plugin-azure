select name, id
from azure.azure_mssql_virtual_machine
where name = 'dummy-test-{{ resourceName }}' and resource_group = '{{ resourceName }}';
