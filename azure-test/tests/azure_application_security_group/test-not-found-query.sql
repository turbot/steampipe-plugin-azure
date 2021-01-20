select name, id, type, region
from azure.azure_application_security_group
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
