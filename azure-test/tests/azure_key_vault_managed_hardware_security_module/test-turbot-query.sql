select name, akas, title, tags
from azure.azure_key_vault_managed_hardware_security_module
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
