select
  id,
  name,
  dns_name,
  region,
  resource_group
from
  azure_bastion_host
where 
  name = 'dummy-{{resourceName}}' 
  and resource_group = '{{resourceName}}'
