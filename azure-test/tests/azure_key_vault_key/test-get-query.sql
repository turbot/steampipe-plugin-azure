select
  name,
  id,
  vault_name,
  enabled,
  key_type,
  curve_name,
  key_size,
  key_uri,
  key_uri_with_version,
  location,
  type,
  key_ops,
  region,
  resource_group,
  subscription_id
from
  azure.azure_key_vault_key
where
  name = '{{ resourceName }}'
  and vault_name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';