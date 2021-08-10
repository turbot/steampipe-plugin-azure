select name, akas, tags, title
from azure.azure_key_vault_managed_hardware_security_module
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}';
