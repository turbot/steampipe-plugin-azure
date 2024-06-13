select
  name,
  id
from
  azure.azure_data_protection_backup_vault
where
  name = '{{ output.resource_name.value }}';