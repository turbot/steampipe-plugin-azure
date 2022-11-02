select name, id, type, platform_update_domain_count, platform_fault_domain_count, region, resource_group, subscription_id
from azure.azure_compute_availability_set
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'