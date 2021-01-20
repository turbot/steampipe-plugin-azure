select name, id, type, active_key_source_vault_id, active_key_url, encryption_tpe, identity_tenant_id, identity_type, region, resource_group, subscription_id
from azure.azure_compute_disk_encryption_set
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'