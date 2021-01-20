select name, akas, tags, title
from azure.azure_application_security_group
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
