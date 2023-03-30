select name, akas, title, tags
from azure.azure_firewall_policy
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';
