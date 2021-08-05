select name, id, region, type, hsm_uri, resource_group, enable_soft_delete, soft_delete_retention_in_days, tenant_id
from azure.azure_key_vault_managed_hardware_security_module
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
