select name, id
from azure.azure_app_service_function_app
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
