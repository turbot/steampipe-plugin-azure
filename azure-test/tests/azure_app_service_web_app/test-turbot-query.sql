select name, akas, title, tags
from azure.azure_app_service_web_app
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';