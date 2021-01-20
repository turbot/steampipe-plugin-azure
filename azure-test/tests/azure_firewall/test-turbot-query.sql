select name, akas, title, tags
from azure.azure_firewall
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
