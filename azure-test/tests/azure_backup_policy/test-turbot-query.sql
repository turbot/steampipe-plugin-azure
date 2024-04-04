select
  name,
  title,
  akas
from
  azure.azure_backup_policy
where
  name = '{{ output.resource_name.value }}'
  and resource_group = '{{ output.resource_name.value }}';