select name, akas, tags, title
from azure.azure_app_service_environment
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
