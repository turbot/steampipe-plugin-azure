select
  lower(id) as id,
	name,
	dns_name,
	resource_group,
	tags
from
	azure_bastion_host
where 
  name = '{{resourceName}}';