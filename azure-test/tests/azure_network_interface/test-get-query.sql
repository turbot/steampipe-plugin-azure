select name, id, region, type, enable_accelerated_networking, enable_ip_forwarding, resource_group, ip_configurations
from azure.azure_network_interface
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
