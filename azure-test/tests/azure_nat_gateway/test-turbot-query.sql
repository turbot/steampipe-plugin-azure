select name, akas, title, tags
from azure_nat_gateway
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';
