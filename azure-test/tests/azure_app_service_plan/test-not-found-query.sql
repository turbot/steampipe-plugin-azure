select name, id
from azure.azure_app_service_plan
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
