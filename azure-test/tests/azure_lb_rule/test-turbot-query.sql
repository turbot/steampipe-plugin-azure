select name, akas, title
from azure.azure_lb_rule
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
