select name, akas, title, tags
from azure.azure_app_service_function_app
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
