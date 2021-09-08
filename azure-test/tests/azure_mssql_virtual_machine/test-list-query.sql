select name, id, type, resource_group, subscription_id, sql_server_license_type
from azure.azure_mssql_virtual_machine
where id = '{{ output.resource_id.value }}';
