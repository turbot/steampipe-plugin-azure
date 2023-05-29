select name, id, provisioning_state, type, resource_group
from azure.azure_firewall_policy
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';
