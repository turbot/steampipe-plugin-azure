select
  name,
  akas,
  title,
  tags
from
  azure.azure_key_vault_key_version
where
  name = '{{ output.key_version.value }}'
  and resource_group = '{{ resourceName }}'
  and vault_name = '{{ resourceName }}';
