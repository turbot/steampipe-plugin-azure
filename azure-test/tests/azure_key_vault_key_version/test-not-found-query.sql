select
  name,
  akas,
  tags,
  title
from
  azure.azure_key_vault_key_version
where
  name = 'dummy-{{ resourceName }}'
  and resource_group = '{{ resourceName }}'
  and vault_name = '{{ resourceName }}';
