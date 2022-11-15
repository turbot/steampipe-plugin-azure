select name, akas, tags, title
from azure_nat_gateway
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}';
