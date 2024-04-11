select
  name,
  id,
  vault_name
from
  azure.azure_backup_policy
where
  name = '{{ output.resource_name.value }}';