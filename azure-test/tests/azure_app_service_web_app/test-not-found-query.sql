select name, id
from azure.azure_app_service_web_app
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
