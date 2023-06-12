select
  id,
  name
from
  azure.azure_key_vault_key_version
where
  name = '{{ output.key_version.value }}';
