select name, akas, title, tags
from azure.azure_app_service_plan
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
