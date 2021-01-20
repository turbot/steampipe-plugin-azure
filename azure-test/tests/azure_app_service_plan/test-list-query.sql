select id, name
from azure.azure_app_service_plan
where name = '{{resourceName}}'
