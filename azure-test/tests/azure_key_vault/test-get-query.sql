select name, id, region, type, vault_uri, resource_group, enabled_for_deployment, enabled_for_disk_encryption, enabled_for_template_deployment, sku_name, tenant_id
from azure.azure_key_vault
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
