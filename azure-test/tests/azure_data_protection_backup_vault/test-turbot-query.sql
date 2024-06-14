select
  name,
  title,
  akas
from
  azure.azure_data_protection_backup_vault
where
  name = '{{ output.resource_name.value }}'
  and resource_group = '{{ output.resource_name.value }}';