select name, id, region
from azure.azure_private_endpoint
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
