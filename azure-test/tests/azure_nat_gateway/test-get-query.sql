select name, id
from azure_nat_gateway
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';
