select name, akas, tags, title
from azure.azure_firewall
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
