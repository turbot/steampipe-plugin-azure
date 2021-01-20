select name, internal_load_balancing_mode, akas, tags, title
from azure.azure_app_service_environment
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
