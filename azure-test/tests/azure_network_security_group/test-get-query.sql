select name, id, type, region, resource_group, security_rules
from azure.azure_network_security_group
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
