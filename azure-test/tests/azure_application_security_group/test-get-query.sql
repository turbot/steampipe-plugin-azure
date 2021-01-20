select name, id, type, region, resource_group
from azure.azure_application_security_group
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
