select name, id, region, front_end_scale_factor, vnet_name, vnet_subnet_name, vnet_resource_group_name, resource_group
from azure.azure_app_service_environment
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
