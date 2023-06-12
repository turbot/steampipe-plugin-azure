select name, akas, title
from azure.azure_firewall_policy
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}';
