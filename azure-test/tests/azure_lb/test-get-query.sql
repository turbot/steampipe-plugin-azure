select name, id, type, region, resource_group, frontend_ip_configurations
from azure.azure_lb
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
