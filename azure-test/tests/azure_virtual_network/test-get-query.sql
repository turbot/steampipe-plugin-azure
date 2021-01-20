select name, id, region, type, enable_ddos_protection, enable_vm_protection, resource_group, address_prefixes
from azure.azure_virtual_network
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
