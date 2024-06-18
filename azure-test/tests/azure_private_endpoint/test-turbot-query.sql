select name, akas, title
from azure.azure_private_endpoint
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
