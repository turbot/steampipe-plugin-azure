select id, name
from azure.azure_app_service_web_app
where name = '{{resourceName}}'
