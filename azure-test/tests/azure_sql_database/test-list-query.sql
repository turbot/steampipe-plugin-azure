select id, name, location, server_name, resource_group
from azure.azure_sql_database
where akas::text = '["{{output.resource_aka.value}}", "{{output.resource_aka_lower.value}}"]';