select id, name
from azure.azure_app_service_environment
where name = '{{resourceName}}'
